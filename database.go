package main

import("fmt"; "encoding/json"; "os"; "sync"; "errors")

type DB struct {
	path string
	mu  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]*Chirp `json:"chirps"`
	Users  map[int]*User  `json:"users"`
}

func NewDB(path string) (*DB, error) {
	db:= &DB{
		path: path,
		mu: &sync.RWMutex{},
	}

	err := db.ensureDB()
	return db, err
}

func (db *DB) writeDB(dbStruct *DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	data, err := json.Marshal(dbStruct)
	if err !=nil {return err}

	return os.WriteFile(db.path, data, 0644)
}

func (db *DB) createDB() error {
	dbStruct := &DBStructure{
		Chirps: map[int]*Chirp{},
		Users:  map[int]*User{},
	}
	return db.writeDB(dbStruct)
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) loadDB() (*DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	data, err := os.ReadFile(db.path)
	if err != nil {return nil, err}

	dbStruct := &DBStructure{}
	err = json.Unmarshal(data, &dbStruct)
	if err != nil {return nil, err}
	return dbStruct, nil
}

func (db *DB) GetChirp(ID int) (*Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {return nil, err}

	chirp, ok := dbStruct.Chirps[ID]
	if !ok {return nil, fmt.Errorf("Error: Chirp doesn't exist")}

	return chirp, nil
}