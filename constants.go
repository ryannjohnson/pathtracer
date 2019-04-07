package pathtracer

// EPS is epsilon, the smallest value this program accepts before it's
// considered "close enough" to zero to be declared zero.
//
// https://github.com/fogleman/pt/blob/999b8a034646a09ad53158cd908c1c01b262aa4d/pt/common.go
// https://stackoverflow.com/questions/626924/what-does-eps-mean-in-c
const EPS = 1e-9
