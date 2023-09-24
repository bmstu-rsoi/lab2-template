package libraries

type Book struct {
	ID        string
	Name      string
	Author    string
	Genre     string
	Condition string
	Available uint64
}

type LibraryBooks struct {
	Total uint64
	Books []Book
}
