package collections

import (
	"fmt"
	"math/rand"
	"reflect"
	"sync"
	"testing"
	"time"
)

func ExampleStack() {
	stack := NewStack[int](1, 2, 3)
	stack.Push(4)
	stack.Push(5)
	stack.Push(6)
	for !stack.IsEmpty() {
		item, _ := stack.Pop()
		fmt.Println(item)
	}
	// Output:
	// 6
	// 5
	// 4
	// 3
	// 2
	// 1
}

func ExampleStack_Peek() {
	stack := NewStack[int](1, 2, 3)
	item, _ := stack.Peek()
	fmt.Println(item)
	// Output: 3
}

func ExampleStack_Pop() {
	stack := NewStack[int](1, 2, 3)
	item, _ := stack.Pop()
	fmt.Println(item)
	// Output: 3
}

func ExampleStack_Push() {
	stack := NewStack[int](1, 2, 3)
	stack.Push(4)
	stack.Push(5)
	stack.Push(6)
	for !stack.IsEmpty() {
		item, _ := stack.Pop()
		fmt.Println(item)
	}
	// Output:
	// 6
	// 5
	// 4
	// 3
	// 2
	// 1
}

func ExampleStack_Size() {
	stack := NewStack[int](1, 2, 3)
	fmt.Println(stack.Size())
	// Output: 3
}

func ExampleStack_String() {
	stack := NewStack[int](1, 2, 3)
	fmt.Println(stack)
	// Output: [1 2 3]
}

func TestNewStack(t *testing.T) {
	type args[T any] struct {
		values []T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want *Stack[T]
	}
	tests := []testCase[string]{
		{
			name: "InitialItems",
			args: args[string]{values: []string{"1", "2", "3"}},
			want: &Stack[string]{items: []string{"1", "2", "3"}},
		},
		{
			name: "NoInitialItems",
			args: args[string]{values: []string{}},
			want: &Stack[string]{items: []string{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStack(tt.args.values...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStack_IsEmpty(t *testing.T) {
	type testCase[T any] struct {
		name string
		s    Stack[T]
		want bool
	}
	tests := []testCase[string]{
		{
			name: "NotEmpty",
			s:    Stack[string]{items: []string{"1", "2", "3"}},
			want: false,
		},
		{
			name: "Empty",
			s:    Stack[string]{items: []string{}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStack_Peek(t *testing.T) {
	type testCase[T any] struct {
		name    string
		s       Stack[T]
		want    T
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name:    "Empty",
			s:       Stack[string]{items: []string{}},
			want:    "",
			wantErr: true,
		},
		{
			name:    "NotEmpty",
			s:       Stack[string]{items: []string{"1", "2", "3"}},
			want:    "3",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Peek()
			if (err != nil) != tt.wantErr {
				t.Errorf("Peek() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Peek() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStack_Pop(t *testing.T) {
	type testCase[T any] struct {
		name    string
		s       Stack[T]
		want    T
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name:    "Empty",
			s:       Stack[string]{items: []string{}},
			want:    "",
			wantErr: true,
		},
		{
			name:    "NotEmpty",
			s:       Stack[string]{items: []string{"1", "2", "3"}},
			want:    "3",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Pop()
			if (err != nil) != tt.wantErr {
				t.Errorf("Pop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pop() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStack_Push(t *testing.T) {
	type args[T any] struct {
		item T
	}
	type testCase[T any] struct {
		name string
		s    Stack[T]
		args args[T]
	}
	tests := []testCase[string]{
		{
			name: "Empty",
			s:    Stack[string]{items: []string{}},
			args: args[string]{item: "1"},
		},
		{
			name: "NotEmpty",
			s:    Stack[string]{items: []string{"1", "2", "3"}},
			args: args[string]{item: "4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Push(tt.args.item)
		})
	}
}

func TestStack_Size(t *testing.T) {
	type testCase[T any] struct {
		name string
		s    Stack[T]
		want int
	}
	tests := []testCase[string]{
		{
			name: "Empty",
			s:    Stack[string]{items: []string{}},
			want: 0,
		},
		{
			name: "NotEmpty",
			s:    Stack[string]{items: []string{"1", "2", "3"}},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Size(); got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStack_String(t *testing.T) {
	type testCase[T any] struct {
		name string
		s    Stack[T]
		want string
	}
	tests := []testCase[string]{
		{
			name: "Empty",
			s:    Stack[string]{items: []string{}},
			want: "[]",
		},
		{
			name: "NotEmpty",
			s:    Stack[string]{items: []string{"1", "2", "3"}},
			want: "[1 2 3]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkStack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		stack := NewStack[int](1, 2, 3)
		stack.Push(4)
		stack.Push(5)
		stack.Push(6)
		for !stack.IsEmpty() {
			_, _ = stack.Pop()
		}
	}
}

func ExampleConcurrentStack() {
	stack := NewConcurrentStack[int]()

	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine for pushing items to the stack
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			stack.Push(i)
			fmt.Printf("Pushed %d\n", i)
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)+1))
		}
	}()

	// Goroutine for popping items from the stack
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			if item, err := stack.Pop(); err == nil {
				fmt.Printf("Popped %d\n", item)
			}
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)+1))
		}
	}()

	wg.Wait()
}

func TestConcurrentStack(t *testing.T) {
	var wg sync.WaitGroup
	stack := NewConcurrentStack[int]() // Replace with your stack implementation

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// sleep between 1ms and 10ms
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)+1))

			// Randomize operations for each goroutine
			switch i % 4 {
			case 0:
				t.Logf("Pushing %d", i)
				stack.Push(i)
			case 1:
				item, err := stack.Pop()
				t.Logf("Popping %d", item)
				if err != nil {
					// Handle error
				}
			case 2:
				_ = stack.IsEmpty()
			case 3:
				_, err := stack.Peek()
				if err != nil {
					// Handle error
				}
			}
		}(i)
	}

	wg.Wait()

	// Perform final state checks and assertions
	// Example: Check if stack is empty or has the expected number of elements
}

func TestConcurrentStack_IsEmpty(t *testing.T) {
	type testCase[T any] struct {
		name string
		s    *ConcurrentStack[T]
		want bool
	}
	tests := []testCase[string]{
		{
			name: "NotEmpty",
			s:    &ConcurrentStack[string]{items: []string{"1", "2", "3"}},
			want: false,
		},
		{
			name: "Empty",
			s:    &ConcurrentStack[string]{items: []string{}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcurrentStack_Peek(t *testing.T) {
	type testCase[T any] struct {
		name    string
		s       *ConcurrentStack[T]
		want    T
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name:    "Empty",
			s:       &ConcurrentStack[string]{items: []string{}},
			want:    "",
			wantErr: true,
		},
		{
			name:    "NotEmpty",
			s:       &ConcurrentStack[string]{items: []string{"1", "2", "3"}},
			want:    "3",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Peek()
			if (err != nil) != tt.wantErr {
				t.Errorf("Peek() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Peek() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcurrentStack_Pop(t *testing.T) {
	type testCase[T any] struct {
		name    string
		s       *ConcurrentStack[T]
		want    T
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name:    "Empty",
			s:       &ConcurrentStack[string]{items: []string{}},
			want:    "",
			wantErr: true,
		},
		{
			name:    "NotEmpty",
			s:       &ConcurrentStack[string]{items: []string{"1", "2", "3"}},
			want:    "3",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Pop()
			if (err != nil) != tt.wantErr {
				t.Errorf("Pop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pop() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcurrentStack_Push(t *testing.T) {
	type args[T any] struct {
		item T
	}
	type testCase[T any] struct {
		name string
		s    *ConcurrentStack[T]
		args args[T]
	}
	tests := []testCase[string]{
		{
			name: "Empty",
			s:    &ConcurrentStack[string]{items: []string{}},
			args: args[string]{item: "1"},
		},
		{
			name: "NotEmpty",
			s:    &ConcurrentStack[string]{items: []string{"1", "2", "3"}},
			args: args[string]{item: "4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Push(tt.args.item)
			if !reflect.DeepEqual(tt.s.items[len(tt.s.items)-1], tt.args.item) {
				t.Errorf("Push() got = %v, want %v", tt.s.items[len(tt.s.items)-1], tt.args.item)
			}
		})
	}
}

func TestConcurrentStack_Size(t *testing.T) {
	type testCase[T any] struct {
		name string
		s    *ConcurrentStack[T]
		want int
	}
	tests := []testCase[string]{
		{
			name: "Empty",
			s:    &ConcurrentStack[string]{items: []string{}},
			want: 0,
		},
		{
			name: "NotEmpty",
			s:    &ConcurrentStack[string]{items: []string{"1", "2", "3"}},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Size(); got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcurrentStack_String(t *testing.T) {
	type testCase[T any] struct {
		name string
		s    *ConcurrentStack[T]
		want string
	}
	tests := []testCase[string]{
		{
			name: "Empty",
			s:    &ConcurrentStack[string]{items: []string{}},
			want: "[]",
		},
		{
			name: "NotEmpty",
			s:    &ConcurrentStack[string]{items: []string{"1", "2", "3"}},
			want: "[1 2 3]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewConcurrentStack(t *testing.T) {
	type args[T any] struct {
		values []T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want *ConcurrentStack[T]
	}
	tests := []testCase[string]{
		{
			name: "InitialItems",
			args: args[string]{values: []string{"1", "2", "3"}},
			want: &ConcurrentStack[string]{items: []string{"1", "2", "3"}},
		},
		{
			name: "NoInitialItems",
			args: args[string]{values: []string{}},
			want: &ConcurrentStack[string]{items: []string{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConcurrentStack(tt.args.values...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConcurrentStack() = %v, want %v", got, tt.want)
			}
		})
	}
}
