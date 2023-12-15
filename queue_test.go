package collections

import (
	"reflect"
	"testing"
)

func TestConcurrentQueue_Clear(t *testing.T) {
	type testCase[T any] struct {
		name string
		q    *ConcurrentQueue[T]
	}
	tests := []testCase[string]{
		{
			name: "Clear",
			q:    NewConcurrentQueue[string]("A", "B", "C"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.q.items) == 0 {
				t.Error("Setup failed")
			}

			tt.q.Clear()

			if len(tt.q.items) != 0 {
				t.Errorf("Clear() got = %v, want %v", tt.q.items, []string{})
			}
		})
	}
}

func TestConcurrentQueue_CopyTo(t *testing.T) {
	type args[T any] struct {
		items []T
		index int
	}
	type testCase[T any] struct {
		name    string
		q       *ConcurrentQueue[T]
		args    args[T]
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name:    "Normal",
			q:       NewConcurrentQueue[string]("A", "B", "C"),
			args:    args[string]{items: make([]string, 3), index: 0},
			wantErr: false,
		},
		{
			name:    "IndexOutOfRange",
			q:       NewConcurrentQueue[string]("A", "B", "C"),
			args:    args[string]{items: make([]string, 3), index: 4},
			wantErr: true,
		},
		{
			name:    "NilArray",
			q:       NewConcurrentQueue[string]("A", "B", "C"),
			args:    args[string]{items: nil, index: 0},
			wantErr: true,
		},
		{
			name:    "DestinationTooSmall",
			q:       NewConcurrentQueue[string]("A", "B", "C"),
			args:    args[string]{items: make([]string, 1), index: 0},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.q.CopyTo(tt.args.items, tt.args.index); (err != nil) != tt.wantErr {
				t.Errorf("CopyTo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConcurrentQueue_Dequeue(t *testing.T) {
	type testCase[T any] struct {
		name    string
		q       *ConcurrentQueue[T]
		want    T
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name:    "Normal",
			q:       NewConcurrentQueue[string]("A", "B", "C"),
			want:    "A",
			wantErr: false,
		},
		{
			name:    "Empty",
			q:       NewConcurrentQueue[string](),
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.q.Dequeue()
			if (err != nil) != tt.wantErr {
				t.Errorf("Dequeue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Dequeue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcurrentQueue_Enqueue(t *testing.T) {
	type args[T any] struct {
		item T
	}
	type testCase[T any] struct {
		name  string
		q     *ConcurrentQueue[T]
		wants []T
		args  args[T]
	}
	tests := []testCase[string]{
		{
			name:  "Normal",
			q:     NewConcurrentQueue[string]("A", "B", "C"),
			wants: []string{"A", "B", "C", "D"},
			args:  args[string]{item: "D"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.Enqueue(tt.args.item)

			if len(tt.q.items) != len(tt.wants) {
				t.Errorf("Enqueue() got = %v, want %v", tt.q.items, tt.wants)
			}
		})
	}
}

func TestConcurrentQueue_IsEmpty(t *testing.T) {
	type testCase[T any] struct {
		name string
		q    *ConcurrentQueue[T]
		want bool
	}
	tests := []testCase[string]{
		{
			name: "Empty",
			q:    NewConcurrentQueue[string](),
			want: true,
		},
		{
			name: "NotEmpty",
			q:    NewConcurrentQueue[string]("A", "B", "C"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcurrentQueue_Peek(t *testing.T) {
	type testCase[T any] struct {
		name    string
		q       *ConcurrentQueue[T]
		want    T
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name:    "Normal",
			q:       NewConcurrentQueue[string]("A", "B", "C"),
			want:    "A",
			wantErr: false,
		},
		{
			name:    "Empty",
			q:       NewConcurrentQueue[string](),
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.q.Peek()
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

func TestConcurrentQueue_Size(t *testing.T) {
	type testCase[T any] struct {
		name string
		q    *ConcurrentQueue[T]
		want int
	}
	tests := []testCase[string]{
		{
			name: "Empty",
			q:    NewConcurrentQueue[string](),
			want: 0,
		},
		{
			name: "NotEmpty",
			q:    NewConcurrentQueue[string]("A", "B", "C"),
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.Size(); got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcurrentQueue_String(t *testing.T) {
	type testCase[T any] struct {
		name string
		q    *ConcurrentQueue[T]
		want string
	}
	tests := []testCase[string]{
		{
			name: "Empty",
			q:    NewConcurrentQueue[string](),
			want: "[]",
		},
		{
			name: "NotEmpty",
			q:    NewConcurrentQueue[string]("A", "B", "C"),
			want: "[A B C]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewConcurrentQueue(t *testing.T) {
	type args[T any] struct {
		values []T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want *ConcurrentQueue[T]
	}
	tests := []testCase[string]{
		{
			name: "Normal",
			args: args[string]{values: []string{"A", "B", "C"}},
			want: &ConcurrentQueue[string]{items: []string{"A", "B", "C"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConcurrentQueue(tt.args.values...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConcurrentQueue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewQueue(t *testing.T) {
	type args[T any] struct {
		values []T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want *Queue[T]
	}
	tests := []testCase[string]{
		{
			name: "Normal",
			args: args[string]{values: []string{"A", "B", "C"}},
			want: &Queue[string]{items: []string{"A", "B", "C"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQueue(tt.args.values...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueue_Clear(t *testing.T) {
	type testCase[T any] struct {
		name string
		q    *Queue[T]
	}
	tests := []testCase[string]{
		{
			name: "Clear",
			q:    NewQueue[string]("A", "B", "C"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.q.items) == 0 {
				t.Error("Setup failed")
			}
			tt.q.Clear()

			if len(tt.q.items) != 0 {
				t.Errorf("Clear() got = %v, want %v", tt.q.items, []string{})
			}
		})
	}
}

func TestQueue_CopyTo(t *testing.T) {
	type args[T any] struct {
		items []T
		index int
	}
	type testCase[T any] struct {
		name    string
		q       *Queue[T]
		args    args[T]
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name:    "Normal",
			q:       NewQueue[string]("A", "B", "C"),
			args:    args[string]{items: make([]string, 3), index: 0},
			wantErr: false,
		},
		{
			name:    "IndexOutOfRange",
			q:       NewQueue[string]("A", "B", "C"),
			args:    args[string]{items: make([]string, 3), index: 4},
			wantErr: true,
		},
		{
			name:    "NilArray",
			q:       NewQueue[string]("A", "B", "C"),
			args:    args[string]{items: nil, index: 0},
			wantErr: true,
		},
		{
			name:    "DestinationTooSmall",
			q:       NewQueue[string]("A", "B", "C"),
			args:    args[string]{items: make([]string, 1), index: 0},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.q.CopyTo(tt.args.items, tt.args.index); (err != nil) != tt.wantErr {
				t.Errorf("CopyTo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQueue_Dequeue(t *testing.T) {
	type testCase[T any] struct {
		name    string
		q       *Queue[T]
		want    T
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name:    "Normal",
			q:       NewQueue[string]("A", "B", "C"),
			want:    "A",
			wantErr: false,
		},
		{
			name:    "Empty",
			q:       NewQueue[string](),
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.q.Dequeue()
			if (err != nil) != tt.wantErr {
				t.Errorf("Dequeue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dequeue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueue_Enqueue(t *testing.T) {
	type args[T any] struct {
		item T
	}
	type testCase[T any] struct {
		name  string
		q     *Queue[T]
		wants []T
		args  args[T]
	}
	tests := []testCase[string]{
		{
			name:  "Normal",
			q:     NewQueue[string]("A", "B", "C"),
			args:  args[string]{item: "D"},
			wants: []string{"A", "B", "C", "D"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.Enqueue(tt.args.item)

			if !reflect.DeepEqual(tt.q.items, tt.wants) {
				t.Errorf("Enqueue() got = %v, want %v", tt.q.items, tt.wants)
			}
		})
	}
}

func TestQueue_IsEmpty(t *testing.T) {
	type testCase[T any] struct {
		name string
		q    *Queue[T]
		want bool
	}
	tests := []testCase[string]{
		{
			name: "Empty",
			q:    NewQueue[string](),
			want: true,
		},
		{
			name: "NotEmpty",
			q:    NewQueue[string]("A", "B", "C"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueue_Peek(t *testing.T) {
	type testCase[T any] struct {
		name    string
		q       *Queue[T]
		want    T
		wantErr bool
	}
	tests := []testCase[string]{
		{
			name:    "Normal",
			q:       NewQueue[string]("A", "B", "C"),
			want:    "A",
			wantErr: false,
		},
		{
			name:    "Empty",
			q:       NewQueue[string](),
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.q.Peek()
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

func TestQueue_Size(t *testing.T) {
	type testCase[T any] struct {
		name string
		q    *Queue[T]
		want int
	}
	tests := []testCase[string]{
		{
			name: "Empty",
			q:    NewQueue[string](),
			want: 0,
		},
		{
			name: "NotEmpty",
			q:    NewQueue[string]("A", "B", "C"),
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.Size(); got != tt.want {
				t.Errorf("Size() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueue_String(t *testing.T) {
	type testCase[T any] struct {
		name string
		q    Queue[T]
		want string
	}
	tests := []testCase[string]{
		{
			name: "Empty",
			q:    *NewQueue[string](),
			want: "[]",
		},
		{
			name: "NotEmpty",
			q:    *NewQueue[string]("A", "B", "C"),
			want: "[A B C]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkConcurrentQueue_Enqueue(b *testing.B) {
	q := NewConcurrentQueue[int]()
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
}

func BenchmarkConcurrentQueue_Dequeue(b *testing.B) {
	q := NewConcurrentQueue[int]()
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = q.Dequeue()
	}
}

func BenchmarkQueue_Enqueue(b *testing.B) {
	q := NewQueue[int]()
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
}

func BenchmarkQueue_Dequeue(b *testing.B) {
	q := NewQueue[int]()
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = q.Dequeue()
	}
}
