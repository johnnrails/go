package internal

type CreateProductUsecase struct {
	ProductRepository ProductRepository
}

type CreateProductInput struct {
	Name  string
	Price float64
}

func NewCreateProductUsecase(r ProductRepository) *CreateProductUsecase {
	return &CreateProductUsecase{r}
}

func (u *CreateProductUsecase) Execute(input CreateProductInput) error {
	product := NewProduct(input.Name, input.Price)
	return u.ProductRepository.Create(product)
}

type ListProductsUsecase struct {
	ProductRepository ProductRepository
}

type ListProductsOutput struct {
	ID    string  `json:"-"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func NewListProductsUsecase(r ProductRepository) *ListProductsUsecase {
	return &ListProductsUsecase{r}
}

func (u *ListProductsUsecase) Execute() ([]*ListProductsOutput, error) {
	products, err := u.ProductRepository.FindAll()

	if err != nil {
		return nil, err
	}

	var output []*ListProductsOutput

	for _, p := range products {
		output = append(output, &ListProductsOutput{
			ID:    p.ID,
			Name:  p.Name,
			Price: p.Price,
		})
	}

	return output, nil
}
