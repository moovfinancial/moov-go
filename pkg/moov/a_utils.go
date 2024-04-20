package moov

/**

Utility functions to help alleviate some annoyances with the API and Go in general

*/

// Can turn a const or hard coded value into a pointer to itself.
func PtrOf[A interface{}](c A) *A {
	return &c
}
