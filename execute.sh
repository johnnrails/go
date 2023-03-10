#!/bin/bash

read -p "Type a digit or a letter > " character
case $character in
  [[:lower:]] | [[:upper:]]) echo "You typed the letter $character" ;;
  [0-9]) echo "You typed the digit $character" ;;
  *) echo "You did not type a lette or a digit"
esac

selection=
until [ "$selection" = "0" ]; do
    echo "
    PROGRAM MENU
    1 - Display free disk space
    2 - Display free memory

    0 - exit program
    "
    echo -n "Enter selection: "
    read selection
    echo ""
    case $selection in
      1 ) df ;;
       ) free ;;
      0 ) exit ;;
      * ) echo "Please enter 1, 2, or 0"
    esac
done
