package server

import (
	"context"
	"fmt"
	"io"
	"math"
	"math/rand"
	"sync"
	"time"

	routepb "github.com/johnnrails/ddd_go/grpc/gen/go/route/v1"
	userpb "github.com/johnnrails/ddd_go/grpc/gen/go/user/v1"
	wearablepb "github.com/johnnrails/ddd_go/grpc/gen/go/wearable/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type WearableServer struct {
}

func (w *WearableServer) BeatsPerMinute(
	r *wearablepb.BeatsPerMinuteRequest,
	stream wearablepb.WearableService_BeatsPerMinuteServer,
) error {
	for {
		select {
		case <-stream.Context().Done():
			return status.Error(codes.Canceled, "Stream has ended")
		default:
			time.Sleep(1 * time.Second)
			value := 30 + rand.Int31n(80)
			err := stream.Send(&wearablepb.BeatsPerMinuteResponse{
				Value:  uint32(value),
				Minute: uint32(time.Now().Second()),
			})
			if err != nil {
				return status.Error(codes.Canceled, "Stream has ended")
			}
		}
	}
}

func (w *WearableServer) ConsumeBeatsPerMinute(
	stream wearablepb.WearableService_ConsumeBeatsPerMinuteServer,
) error {
	var total int
	for {
		value, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&wearablepb.ConsumeBeatsPerMinuteResponse{
				Total: uint32(total),
			})
		}

		if err != nil {
			return err
		}

		fmt.Println(value.GetMinute(), value.GetValue(), value.GetUuid())
		total++
	}
}

type UserService struct {
	userpb.UserServiceServer
}

func (u *UserService) GetUser(_ context.Context, r *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	return &userpb.GetUserResponse{
		User: &userpb.User{
			Uuid:          r.Uuid,
			Fullname:      "Johnn",
			BirthYear:     2007,
			Salary:        10,
			Addresses:     []*userpb.Address{},
			MaritalStatus: userpb.MaritalStatus_MARITAL_STATUS_MARRIED,
		},
	}, nil
}

type routeGuideServer struct {
	routepb.UnimplementedRouteGuideServer
	savedFeatures []*routepb.Feature // read-only after initialized

	mu         sync.Mutex // protects routeNotes
	routeNotes map[string][]*routepb.RouteNote
}

func (s *routeGuideServer) GetFeature(ctx context.Context, p *routepb.Point) (*routepb.Feature, error) {
	for _, f := range s.savedFeatures {
		if proto.Equal(f.Location, p) {
			return f, nil
		}
	}
	return &routepb.Feature{Location: p}, nil
}

func (s *routeGuideServer) ListFeatures(rect *routepb.Rectangle, stream routepb.RouteGuide_ListFeaturesServer) error {
	for _, f := range s.savedFeatures {
		if inRange(f.Location, rect) {
			if err := stream.Send(f); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *routeGuideServer) RecordRoute(stream routepb.RouteGuide_RecordRouteServer) error {
	var pointsCount, featuresCount, distance int32
	var lastPoint *routepb.Point
	startTime := time.Now()
	for {
		point, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&routepb.RouteSummary{
				PointCount:   pointsCount,
				FeatureCount: featuresCount,
				Distance:     distance,
				ElapsedTime:  int32(time.Now().Sub(startTime).Seconds()),
			})
		}
		if err != nil {
			return err
		}
		pointsCount++
		for _, f := range s.savedFeatures {
			if proto.Equal(f.Location, point) {
				featuresCount++
			}
		}
		if lastPoint != nil {
			distance += calcDistance(lastPoint, point)
		}
		lastPoint = point
	}
}

func inRange(point *routepb.Point, rect *routepb.Rectangle) bool {
	left := math.Min(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
	right := math.Max(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
	top := math.Max(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))
	bottom := math.Min(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))

	if float64(point.Longitude) >= left &&
		float64(point.Longitude) <= right &&
		float64(point.Latitude) >= bottom &&
		float64(point.Latitude) <= top {
		return true
	}
	return false
}

func toRadians(num float64) float64 {
	return num * math.Pi / float64(180)
}

// calcDistance calculates the distance between two points using the "haversine" formula.
// The formula is based on http://mathforum.org/library/drmath/view/51879.html.
func calcDistance(p1 *routepb.Point, p2 *routepb.Point) int32 {
	const CordFactor float64 = 1e7
	const R = float64(6371000) // earth radius in metres
	lat1 := toRadians(float64(p1.Latitude) / CordFactor)
	lat2 := toRadians(float64(p2.Latitude) / CordFactor)
	lng1 := toRadians(float64(p1.Longitude) / CordFactor)
	lng2 := toRadians(float64(p2.Longitude) / CordFactor)
	dlat := lat2 - lat1
	dlng := lng2 - lng1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlng/2)*math.Sin(dlng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c
	return int32(distance)
}
