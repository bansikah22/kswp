package models

import "time"

type Resource struct {
	Name      string
	Namespace string
	Kind      string
	Reason    string
	Age       time.Duration
}
