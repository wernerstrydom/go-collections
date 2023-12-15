package collections

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func TestNewConcurrentList(t *testing.T) {
	type args[T comparable] struct {
		values []T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want *ConcurrentList[T]
	}
	tests := []testCase[string]{
		{
			name: "Normal",
			args: args[string]{values: []string{"1", "2", "3"}},
			want: &ConcurrentList[string]{items: []string{"1", "2", "3"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewConcurrentList(tt.args.values...)
			if !reflect.DeepEqual(got.items, tt.want.items) {
				t.Errorf("NewConcurrentList() = %v, want %v", got, tt.want)
			}
			if got.comparer == nil {
				t.Error("NewConcurrentList() returned a ConcurrentList with a nil comparer")
			}
		})
	}
}

func TestNewConcurrentListWithComparer(t *testing.T) {
	c := func(a, b string) bool { return strings.Compare(strings.ToLower(a), strings.ToLower(b)) == 0 }

	type args[T any] struct {
		comparer EqualityComparer[T]
		values   []T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want *ConcurrentList[T]
	}
	tests := []testCase[string]{
		{
			name: "Normal",
			args: args[string]{comparer: c, values: []string{"1", "2", "3"}},
			want: &ConcurrentList[string]{items: []string{"1", "2", "3"}, comparer: c},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewConcurrentListWithComparer(tt.args.comparer, tt.args.values...)
			if !reflect.DeepEqual(got.items, tt.want.items) {
				t.Errorf("NewConcurrentListWithComparer() = %v, want %v", got, tt.want)
			}
			if got.comparer == nil {
				t.Errorf("NewConcurrentListWithComparer() returned a ConcurrentList with a nil comparer")
			}
		})
	}
}

func TestConcurrentList_Add(t *testing.T) {
	type args[T any] struct {
		item T
	}
	type testCase[T any] struct {
		name string
		l    *ConcurrentList[T]
		args args[T]
	}
	tests := []testCase[string]{
		{
			name: "Default",
			l:    NewConcurrentList[string](),
			args: args[string]{item: "1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Add(tt.args.item)
			found := false
			for _, v := range tt.l.items {
				if v == tt.args.item {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Add() did not add %v to the list", tt.args.item)
			}

		})
	}
}

func TestConcurrentList_AddRange(t *testing.T) {
	type args[T any] struct {
		items []T
	}
	type testCase[T any] struct {
		name string
		l    *ConcurrentList[T]
		args args[T]
	}
	tests := []testCase[string]{
		{
			name: "Default",
			l:    NewConcurrentList[string](),
			args: args[string]{items: []string{"1", "2", "3"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.AddRange(tt.args.items...)

			for _, v := range tt.args.items {
				found := false
				for _, v2 := range tt.l.items {
					if v == v2 {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("AddRange() did not add %v to the list", v)
				}
			}
		})
	}
}

func TestConcurrentList_Clear(t *testing.T) {
	type testCase[T any] struct {
		name string
		l    *ConcurrentList[T]
	}
	tests := []testCase[string]{
		{
			name: "Default",
			l:    NewConcurrentList[string]("1", "2", "3"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if len(tt.l.items) == 0 {
				t.Errorf("Setup wasn't executed")
			}

			tt.l.Clear()

			if len(tt.l.items) != 0 {
				t.Errorf("Clear() did not clear the list")
			}
		})
	}
}

func TestConcurrentList_Contains(t *testing.T) {
	type args[T any] struct {
		item T
	}
	type testCase[T any] struct {
		name string
		l    *ConcurrentList[T]
		args args[T]
		want bool
	}
	tests := []testCase[string]{
		{
			name: "Existing",
			l:    NewConcurrentList[string]("1", "2", "3"),
			args: args[string]{item: "2"},
			want: true,
		},
		{
			name: "NonExisting",
			l:    NewConcurrentList[string]("1", "2", "3"),
			args: args[string]{item: "4"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Contains(tt.args.item); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcurrentList_CopyTo(t *testing.T) {
	type args[T any] struct {
		array      []T
		arrayIndex int
	}
	type testCase[T any] struct {
		name    string
		l       *ConcurrentList[T]
		args    args[T]
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name:    "Normal",
			l:       NewConcurrentList[string]("1", "2", "3"),
			args:    args[string]{array: make([]string, 3), arrayIndex: 0},
			wantErr: false,
		},
		{
			name:    "IndexOutOfRange",
			l:       NewConcurrentList[string]("1", "2", "3"),
			args:    args[string]{array: make([]string, 3), arrayIndex: 4},
			wantErr: true,
		},
		{
			name:    "NilArray",
			l:       NewConcurrentList[string]("1", "2", "3"),
			args:    args[string]{array: nil, arrayIndex: 0},
			wantErr: true,
		},
		{
			name:    "DestinationTooSmall",
			l:       NewConcurrentList[string]("1", "2", "3"),
			args:    args[string]{array: make([]string, 2), arrayIndex: 0},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.l.CopyTo(tt.args.array, tt.args.arrayIndex); (err != nil) != tt.wantErr {
				t.Errorf("CopyTo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConcurrentList_Get(t *testing.T) {
	type args struct {
		index int
	}
	type testCase[T any] struct {
		name    string
		l       *ConcurrentList[T]
		args    args
		want    T
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name:    "Normal",
			l:       NewConcurrentList[string]("1", "2", "3"),
			args:    args{index: 1},
			want:    "2",
			wantErr: false,
		},
		{
			name:    "IndexOutOfRange",
			l:       NewConcurrentList[string]("1", "2", "3"),
			args:    args{index: 4},
			want:    "",
			wantErr: true,
		},
		{
			name:    "NegativeIndex",
			l:       NewConcurrentList[string]("1", "2", "3"),
			args:    args{index: -1},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.l.Get(tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcurrentList_IndexOf(t *testing.T) {
	type args[T any] struct {
		item T
	}
	type testCase[T any] struct {
		name string
		l    *ConcurrentList[T]
		args args[T]
		want int
	}
	tests := []testCase[string]{
		{
			name: "InRange",
			l:    NewConcurrentList[string]("1", "2", "3"),
			args: args[string]{item: "2"},
			want: 1,
		},
		{
			name: "OutOfRange",
			l:    NewConcurrentList[string]("1", "2", "3"),
			args: args[string]{item: "4"},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.IndexOf(tt.args.item); got != tt.want {
				t.Errorf("IndexOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcurrentList_Insert(t *testing.T) {
	type args[T any] struct {
		index int
		item  T
	}
	type testCase[T any] struct {
		name     string
		l        *ConcurrentList[T]
		args     args[T]
		wantList []T
		wantErr  bool
	}
	tests := []testCase[string]{
		{
			name:     "Normal",
			l:        NewConcurrentList[string]("1", "2", "3"),
			args:     args[string]{index: 1, item: "4"},
			wantList: []string{"1", "4", "2", "3"},
			wantErr:  false,
		},
		{
			name:     "IndexOutOfRange",
			l:        NewConcurrentList[string]("1", "2", "3"),
			args:     args[string]{index: 4, item: "4"},
			wantList: []string{"1", "2", "3"},
			wantErr:  true,
		},
		{
			name:     "NegativeIndex",
			l:        NewConcurrentList[string]("1", "2", "3"),
			args:     args[string]{index: -1, item: "4"},
			wantList: []string{"1", "2", "3"},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.l.Insert(tt.args.index, tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(tt.l.items, tt.wantList) {
				t.Errorf("Insert() got = %v, want %v", tt.l.items, tt.wantList)
			}
		})
	}
}

func TestConcurrentList_LastIndexOf(t *testing.T) {
	type args[T any] struct {
		item T
	}
	type testCase[T any] struct {
		name string
		l    *ConcurrentList[T]
		args args[T]

		want int
	}
	tests := []testCase[string]{
		{
			name: "InRange",
			l:    NewConcurrentList[string]("1", "2", "3", "2"),
			args: args[string]{item: "2"},
			want: 3,
		},
		{
			name: "OutOfRange",
			l:    NewConcurrentList[string]("1", "2", "3"),
			args: args[string]{item: "4"},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.LastIndexOf(tt.args.item); got != tt.want {
				t.Errorf("LastIndexOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcurrentList_Remove(t *testing.T) {
	type args[T any] struct {
		item T
	}
	type testCase[T any] struct {
		name     string
		l        *ConcurrentList[T]
		args     args[T]
		wantList []T
		want     bool
	}
	tests := []testCase[string]{
		{
			name:     "Existing",
			l:        NewConcurrentList[string]("1", "2", "3"),
			args:     args[string]{item: "2"},
			wantList: []string{"1", "3"},
			want:     true,
		},
		{
			name:     "NonExisting",
			l:        NewConcurrentList[string]("1", "2", "3"),
			args:     args[string]{item: "4"},
			wantList: []string{"1", "2", "3"},
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Remove(tt.args.item); got != tt.want {
				t.Errorf("Remove() = %v, want %v", got, tt.want)
				return
			}
			if !reflect.DeepEqual(tt.l.items, tt.wantList) {
				t.Errorf("Remove() got = %v, want %v", tt.l.items, tt.wantList)
			}
		})
	}
}

func TestConcurrentList_RemoveAt(t *testing.T) {
	type args struct {
		index int
	}
	type testCase[T any] struct {
		name     string
		l        *ConcurrentList[T]
		args     args
		wantList []T
		wantErr  bool
	}
	tests := []testCase[string]{
		{
			name:     "Normal",
			l:        NewConcurrentList[string]("1", "2", "3"),
			args:     args{index: 1},
			wantList: []string{"1", "3"},
			wantErr:  false,
		},
		{
			name:     "IndexOutOfRange",
			l:        NewConcurrentList[string]("1", "2", "3"),
			args:     args{index: 4},
			wantList: []string{"1", "2", "3"},
			wantErr:  true,
		},
		{
			name:     "NegativeIndex",
			l:        NewConcurrentList[string]("1", "2", "3"),
			args:     args{index: -1},
			wantList: []string{"1", "2", "3"},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.l.RemoveAt(tt.args.index); (err != nil) != tt.wantErr {
				t.Errorf("RemoveAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.l.items, tt.wantList) {
				t.Errorf("RemoveAt() got = %v, want %v", tt.l.items, tt.wantList)
			}
		})
	}
}

func TestConcurrentList_Set(t *testing.T) {
	type args[T any] struct {
		index int
		item  T
	}
	type testCase[T any] struct {
		name     string
		l        *ConcurrentList[T]
		args     args[T]
		wantList []T
		wantErr  bool
	}
	tests := []testCase[string]{
		{
			name:     "Normal",
			l:        NewConcurrentList[string]("1", "2", "3"),
			args:     args[string]{index: 1, item: "4"},
			wantList: []string{"1", "4", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.l.Set(tt.args.index, tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func ExampleConcurrentList() {
	// Create a new ConcurrentList
	list := NewConcurrentList[int]()

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Launch several goroutines that add elements to the list
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			list.Add(i)
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Print the contents of the list
	fmt.Println(list.String())
}

func ExampleConcurrentList_Add() {
	// Create a new ConcurrentList
	list := NewConcurrentList[int]()

	// Add some items to the list
	list.Add(1)
	list.Add(2)
	list.Add(3)

	// Print the contents of the list
	fmt.Println(list.String())
	// Output: [1 2 3]
}

func ExampleConcurrentList_AddRange() {
	// Create a new ConcurrentList
	list := NewConcurrentList[int]()

	// Add some items to the list
	list.AddRange(1, 2, 3)

	// Print the contents of the list
	fmt.Println(list.String())
	// Output: [1 2 3]
}

func ExampleConcurrentList_Clear() {
	// Create a new ConcurrentList
	list := NewConcurrentList[int](1, 2, 3)

	// Clear the list
	list.Clear()

	// Print the contents of the list
	fmt.Println(list.String())
	// Output: []
}

func ExampleConcurrentList_Contains() {
	// Create a new ConcurrentList
	list := NewConcurrentList[int](1, 2, 3)

	// Check if the list contains 2
	fmt.Println(list.Contains(2))
	// Output: true
}

func ExampleConcurrentList_CopyTo() {
	// Create a new ConcurrentList
	list := NewConcurrentList[int](1, 2, 3)

	// Create a destination array
	array := make([]int, 3)

	// Copy the list to the array
	_ = list.CopyTo(array, 0)

	// Print the contents of the array
	fmt.Println(array)
	// Output: [1 2 3]
}

func ExampleConcurrentList_Get() {
	// Create a new ConcurrentList
	list := NewConcurrentList[int](1, 2, 3)

	// Get the item at index 1
	item, _ := list.Get(1)

	// Print the item
	fmt.Println(item)
	// Output: 2
}

func ExampleConcurrentList_IndexOf() {
	// Create a new ConcurrentList
	list := NewConcurrentList[int](1, 2, 3)

	// Get the index of 2
	index := list.IndexOf(2)

	// Print the index
	fmt.Println(index)
	// Output: 1
}

func ExampleConcurrentList_Insert() {
	// Create a new ConcurrentList
	list := NewConcurrentList[int](1, 2, 3)

	// Insert 4 at index 1
	_ = list.Insert(1, 4)

	// Print the contents of the list
	fmt.Println(list.String())
	// Output: [1 4 2 3]
}

func ExampleConcurrentList_LastIndexOf() {
	// Create a new ConcurrentList
	list := NewConcurrentList[int](1, 2, 3, 2)

	// Get the last index of 2
	index := list.LastIndexOf(2)

	// Print the index
	fmt.Println(index)
	// Output: 3
}

func ExampleConcurrentList_Remove() {
	// Create a new ConcurrentList
	list := NewConcurrentList[int](1, 2, 3)

	// Remove 2 from the list
	_ = list.Remove(2)

	// Print the contents of the list
	fmt.Println(list.String())
	// Output: [1 3]
}

func ExampleConcurrentList_RemoveAt() {
	// Create a new ConcurrentList
	list := NewConcurrentList[int](1, 2, 3)

	// Remove the item at index 1
	_ = list.RemoveAt(1)

	// Print the contents of the list
	fmt.Println(list.String())
	// Output: [1 3]
}

func ExampleConcurrentList_Set() {
	// Create a new ConcurrentList
	list := NewConcurrentList[int](1, 2, 3)

	// Set the item at index 1 to 4
	_ = list.Set(1, 4)

	// Print the contents of the list
	fmt.Println(list.String())
	// Output: [1 4 3]
}

func TestList_Add(t *testing.T) {
	type args[T any] struct {
		item T
	}
	type testCase[T any] struct {
		name     string
		l        *List[T]
		args     args[T]
		wantList []T
	}
	tests := []testCase[string]{
		{
			name:     "Normal",
			l:        NewList[string]("1", "2", "3"),
			args:     args[string]{item: "4"},
			wantList: []string{"1", "2", "3", "4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.Add(tt.args.item)
			if !reflect.DeepEqual(tt.l.items, tt.wantList) {
				t.Errorf("Add() got = %v, want %v", tt.l.items, tt.wantList)
			}
		})
	}
}

func TestList_AddRange(t *testing.T) {
	type args[T any] struct {
		items []T
	}
	type testCase[T any] struct {
		name string
		l    *List[T]
		args args[T]
		want []T
	}
	tests := []testCase[string]{
		{
			name: "Normal",
			l:    NewList[string]("1", "2", "3"),
			args: args[string]{items: []string{"4", "5", "6"}},
			want: []string{"1", "2", "3", "4", "5", "6"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.l.AddRange(tt.args.items...)

			if !reflect.DeepEqual(tt.l.items, tt.want) {
				t.Errorf("AddRange() got = %v, want %v", tt.l.items, tt.want)
			}
		})
	}
}

func TestList_Clear(t *testing.T) {
	type testCase[T any] struct {
		name string
		l    *List[T]
	}
	tests := []testCase[string]{
		{
			name: "Normal",
			l:    NewList[string]("1", "2", "3"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.l.items) == 0 {
				t.Errorf("Setup wasn't executed")
			}
			tt.l.Clear()
			if len(tt.l.items) != 0 {
				t.Errorf("Clear() did not clear the list")
			}
		})
	}
}

func TestList_Contains(t *testing.T) {
	type args[T any] struct {
		item T
	}
	type testCase[T any] struct {
		name string
		l    *List[T]
		args args[T]
		want bool
	}
	tests := []testCase[string]{
		{
			name: "Existing",
			l:    NewList[string]("1", "2", "3"),
			args: args[string]{item: "2"},
			want: true,
		},
		{
			name: "NonExisting",
			l:    NewList[string]("1", "2", "3"),
			args: args[string]{item: "4"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Contains(tt.args.item); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList_CopyTo(t *testing.T) {
	type args[T any] struct {
		array      []T
		arrayIndex int
	}
	type testCase[T any] struct {
		name    string
		l       *List[T]
		args    args[T]
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name:    "Normal",
			l:       NewList[string]("1", "2", "3"),
			args:    args[string]{array: make([]string, 3), arrayIndex: 0},
			wantErr: false,
		},
		{
			name:    "IndexOutOfRange",
			l:       NewList[string]("1", "2", "3"),
			args:    args[string]{array: make([]string, 3), arrayIndex: 4},
			wantErr: true,
		},
		{
			name:    "NilArray",
			l:       NewList[string]("1", "2", "3"),
			args:    args[string]{array: nil, arrayIndex: 0},
			wantErr: true,
		},
		{
			name:    "DestinationTooSmall",
			l:       NewList[string]("1", "2", "3"),
			args:    args[string]{array: make([]string, 2), arrayIndex: 0},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.l.CopyTo(tt.args.array, tt.args.arrayIndex); (err != nil) != tt.wantErr {
				t.Errorf("CopyTo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestList_Get(t *testing.T) {
	type args struct {
		index int
	}
	type testCase[T any] struct {
		name    string
		l       *List[T]
		args    args
		want    T
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name:    "Normal",
			l:       NewList[string]("1", "2", "3"),
			args:    args{index: 1},
			want:    "2",
			wantErr: false,
		},
		{
			name:    "IndexOutOfRange",
			l:       NewList[string]("1", "2", "3"),
			args:    args{index: 4},
			want:    "",
			wantErr: true,
		},
		{
			name:    "NegativeIndex",
			l:       NewList[string]("1", "2", "3"),
			args:    args{index: -1},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.l.Get(tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList_IndexOf(t *testing.T) {
	type args[T any] struct {
		item T
	}
	type testCase[T any] struct {
		name string
		l    *List[T]
		args args[T]
		want int
	}
	tests := []testCase[string]{
		{
			name: "InRange",
			l:    NewList[string]("1", "2", "3"),
			args: args[string]{item: "2"},
			want: 1,
		},
		{
			name: "OutOfRange",
			l:    NewList[string]("1", "2", "3"),
			args: args[string]{item: "4"},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.IndexOf(tt.args.item); got != tt.want {
				t.Errorf("IndexOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList_Insert(t *testing.T) {
	type args[T any] struct {
		index int
		item  T
	}
	type testCase[T any] struct {
		name    string
		l       *List[T]
		args    args[T]
		want    []T
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name:    "Normal",
			l:       NewList[string]("1", "2", "3"),
			args:    args[string]{index: 1, item: "4"},
			want:    []string{"1", "4", "2", "3"},
			wantErr: false,
		},
		{
			name:    "IndexOutOfRange",
			l:       NewList[string]("1", "2", "3"),
			args:    args[string]{index: 4, item: "4"},
			want:    []string{"1", "2", "3"},
			wantErr: true,
		},
		{
			name:    "NegativeIndex",
			l:       NewList[string]("1", "2", "3"),
			args:    args[string]{index: -1, item: "4"},
			want:    []string{"1", "2", "3"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.l.Insert(tt.args.index, tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(tt.l.items, tt.want) {
				t.Errorf("Insert() got = %v, want %v", tt.l.items, tt.want)
			}
		})
	}
}

func TestList_LastIndexOf(t *testing.T) {
	type args[T any] struct {
		item T
	}
	type testCase[T any] struct {
		name string
		l    *List[T]
		args args[T]
		want int
	}
	tests := []testCase[string]{
		{
			name: "InRange",
			l:    NewList[string]("1", "2", "3", "2"),
			args: args[string]{item: "2"},
			want: 3,
		},
		{
			name: "OutOfRange",
			l:    NewList[string]("1", "2", "3"),
			args: args[string]{item: "4"},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.LastIndexOf(tt.args.item); got != tt.want {
				t.Errorf("LastIndexOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList_Remove(t *testing.T) {
	type args[T any] struct {
		item T
	}
	type testCase[T any] struct {
		name string
		l    *List[T]
		args args[T]
		want bool
	}
	tests := []testCase[string]{
		{
			name: "Existing",
			l:    NewList[string]("1", "2", "3"),
			args: args[string]{item: "2"},
			want: true,
		},
		{
			name: "NonExisting",
			l:    NewList[string]("1", "2", "3"),
			args: args[string]{item: "4"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Remove(tt.args.item); got != tt.want {
				t.Errorf("Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList_RemoveAt(t *testing.T) {
	type args struct {
		index int
	}
	type testCase[T any] struct {
		name    string
		l       *List[T]
		args    args
		want    []T
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name:    "Normal",
			l:       NewList[string]("1", "2", "3"),
			args:    args{index: 1},
			want:    []string{"1", "3"},
			wantErr: false,
		},
		{
			name:    "IndexOutOfRange",
			l:       NewList[string]("1", "2", "3"),
			args:    args{index: 4},
			want:    []string{"1", "2", "3"},
			wantErr: true,
		},
		{
			name:    "NegativeIndex",
			l:       NewList[string]("1", "2", "3"),
			args:    args{index: -1},
			want:    []string{"1", "2", "3"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.l.RemoveAt(tt.args.index); (err != nil) != tt.wantErr {
				t.Errorf("RemoveAt() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.l.items, tt.want) {
				t.Errorf("RemoveAt() got = %v, want %v", tt.l.items, tt.want)
			}
		})
	}
}

func TestList_Set(t *testing.T) {
	type args[T any] struct {
		index int
		item  T
	}
	type testCase[T any] struct {
		name    string
		l       *List[T]
		args    args[T]
		want    []T
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name:    "Normal",
			l:       NewList[string]("1", "2", "3"),
			args:    args[string]{index: 1, item: "4"},
			want:    []string{"1", "4", "3"},
			wantErr: false,
		},
		{
			name:    "IndexOutOfRange",
			l:       NewList[string]("1", "2", "3"),
			args:    args[string]{index: 4, item: "4"},
			want:    []string{"1", "2", "3"},
			wantErr: true,
		},
		{
			name:    "NegativeIndex",
			l:       NewList[string]("1", "2", "3"),
			args:    args[string]{index: -1, item: "4"},
			want:    []string{"1", "2", "3"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.l.Set(tt.args.index, tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(tt.l.items, tt.want) {
				t.Errorf("Set() got = %v, want %v", tt.l.items, tt.want)
			}
		})
	}
}

func TestNewList(t *testing.T) {
	type args[T comparable] struct {
		values []T
	}
	type testCase[T comparable] struct {
		name string
		args args[T]
		want *List[T]
	}
	tests := []testCase[string]{
		{
			name: "Normal",
			args: args[string]{values: []string{"1", "2", "3"}},
			want: &List[string]{items: []string{"1", "2", "3"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewList(tt.args.values...)

			if !reflect.DeepEqual(got.items, tt.want.items) {
				t.Errorf("NewList() = %v, want %v", got, tt.want)
			}

			if got.comparer == nil {
				t.Errorf("NewList() returned a List with a nil comparer")
			}
		})
	}
}

func TestNewListWithEqualityComparer(t *testing.T) {
	type args[T any] struct {
		comparer EqualityComparer[T]
		values   []T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want *List[T]
	}
	comparer := func(a, b string) bool { return strings.ToLower(a) == strings.ToLower(b) }
	tests := []testCase[string]{
		{
			name: "Normal",
			args: args[string]{comparer: comparer, values: []string{"1", "2", "3"}},
			want: &List[string]{items: []string{"1", "2", "3"}, comparer: comparer},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewListWithEqualityComparer(tt.args.comparer, tt.args.values...)
			if !reflect.DeepEqual(got.items, tt.want.items) {
				t.Errorf("NewListWithEqualityComparer() = %v, want %v", got, tt.want)
			}
			if got.comparer == nil {
				t.Errorf("NewListWithEqualityComparer() returned a List with a nil comparer")
			}
		})
	}
}

func TestList_String(t *testing.T) {
	type testCase[T any] struct {
		name string
		l    *List[T]
		want string
	}
	tests := []testCase[string]{
		{
			name: "Normal",
			l:    NewList[string]("1", "2", "3"),
			want: "[1 2 3]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleNewList() {
	// Create a new List
	list := NewList[int](1, 2, 3)

	// Print the contents of the list
	fmt.Println(list.String())
	// Output: [1 2 3]
}

func ExampleNewListWithEqualityComparer() {
	// Create a new List with a case-insensitive comparer
	list := NewListWithEqualityComparer[string](func(a, b string) bool {
		return strings.ToLower(a) == strings.ToLower(b)
	}, "A", "b", "C")

	// Print the contents of the list
	fmt.Println(list.String())

	list.Remove("A")

	// Print the contents of the list
	fmt.Println(list.String())

	// Output: [A b C]
	// [b C]
}

func ExampleList_Add() {
	// Create a new List
	list := NewList[int](1, 2, 3)

	// Add an item to the list
	list.Add(4)

	// Print the contents of the list
	fmt.Println(list.String())
	// Output: [1 2 3 4]
}

func ExampleList_AddRange() {
	// Create a new List
	list := NewList[int](1, 2, 3)

	// Add some items to the list
	list.AddRange(4, 5, 6)

	// Print the contents of the list
	fmt.Println(list.String())
	// Output: [1 2 3 4 5 6]
}

func ExampleList_Clear() {
	// Create a new List
	list := NewList[int](1, 2, 3)

	// Clear the list
	list.Clear()

	// Print the contents of the list
	fmt.Println(list.String())
	// Output: []
}

func ExampleList_Contains() {
	// Create a new List
	list := NewList[int](1, 2, 3)

	// Check if the list contains 2
	fmt.Println(list.Contains(2))
	// Output: true
}

func ExampleList_CopyTo() {
	// Create a new List
	list := NewList[int](1, 2, 3)

	// Create a destination array
	array := make([]int, 3)

	// Copy the list to the array
	_ = list.CopyTo(array, 0)

	// Print the contents of the array
	fmt.Println(array)
	// Output: [1 2 3]
}

func ExampleList_Get() {
	// Create a new List
	list := NewList[int](1, 2, 3)

	// Get the item at index 1
	item, _ := list.Get(1)

	// Print the item
	fmt.Println(item)
	// Output: 2
}

func ExampleList_IndexOf() {
	// Create a new List
	list := NewList[int](1, 2, 3)

	// Get the index of 2
	index := list.IndexOf(2)

	// Print the index
	fmt.Println(index)
	// Output: 1
}

func ExampleList_Insert() {
	// Create a new List
	list := NewList[int](1, 2, 3)

	// Insert 4 at index 1
	_ = list.Insert(1, 4)

	// Print the contents of the list
	fmt.Println(list.String())
	// Output: [1 4 2 3]
}

func ExampleList_LastIndexOf() {
	// Create a new List
	list := NewList[int](1, 2, 3, 2)

	// Get the last index of 2
	index := list.LastIndexOf(2)

	// Print the index
	fmt.Println(index)
	// Output: 3
}

func ExampleList_Remove() {
	// Create a new List
	list := NewList[int](1, 2, 3)

	// Remove 2 from the list
	_ = list.Remove(2)

	// Print the contents of the list
	fmt.Println(list.String())
	// Output: [1 3]
}

func ExampleList_RemoveAt() {
	// Create a new List
	list := NewList[int](1, 2, 3)

	// Remove the item at index 1
	_ = list.RemoveAt(1)

	// Print the contents of the list
	fmt.Println(list.String())
	// Output: [1 3]
}

func ExampleList_Set() {
	// Create a new List
	list := NewList[int](1, 2, 3)

	// Set the item at index 1 to 4
	_ = list.Set(1, 4)

	// Print the contents of the list
	fmt.Println(list.String())
	// Output: [1 4 3]
}
