package snapshot

import (
	"encoding/json"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/recruit-tech/RISUCON2022Summer/snapshots/generator/model"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var snapshots []model.Snapshot

func LoadSnapshots(snapshotDir string) error {
	paths, err := filepath.Glob(filepath.Join(snapshotDir, "snapshot*.json"))
	if err != nil {
		return err
	}

	path := paths[rand.Intn(len(paths))]

	w, err := os.Open(path)
	if err != nil {
		return err
	}
	defer w.Close()

	err = json.NewDecoder(w).Decode(&snapshots)
	if err != nil {
		return err
	}

	return nil
}

func GetSnapshots() []model.Snapshot {
	return snapshots
}

var (
	idMap map[string]string = map[string]string{
		"000000000006AFVGQT5ZYC0GEK": "000000000006AFVGQT5ZYC0GEK",
		"0000000001KFV89NC6TD3MX7K4": "0000000001KFV89NC6TD3MX7K4",
		"00000000026XVYZFKH7QNBCC19": "00000000026XVYZFKH7QNBCC19",
	}
	mu sync.RWMutex
)

func SetID(key, id string) {
	mu.Lock()
	defer mu.Unlock()

	idMap[key] = id
}

func GetID(key string) string {
	mu.RLock()
	defer mu.RUnlock()

	return idMap[key]
}
