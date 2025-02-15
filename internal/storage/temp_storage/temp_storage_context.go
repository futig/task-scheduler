package tempstorage

import "sync"

type TempStorageContext struct {
	Tasks   sync.Map
	Users   sync.Map
	Reminds sync.Map
}
