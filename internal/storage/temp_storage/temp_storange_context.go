package tempstorage

import "sync"

type TempStorageContext struct {
	Tasks   sync.Map
	Reminds sync.Map
}
