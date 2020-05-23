package pgkebab

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

// DBLink is the basic interface between developer and database instance
type DBLink struct {
	db            *sql.DB // db DBLink Connection Handler
	DBLinkOptions         // connection options
	supposedReady bool
}

// log prints log messages in case of l.debugPrint is true
func (l *DBLink) log(format string, args ...interface{}) {
	if l.debugPrint {
		s := fmt.Sprintf(format, args...)

		log.Println(s)
	}
}

// SetEmergencyCallback defines a routine to log emergency "no-connection" messages
func (l *DBLink) SetEmergencyCallback(ec func(error)) {
	l.emergencyCallback = ec
}

func (l *DBLink) connectSimple() error {
	l.supposedReady = false
	l.db = nil

	db, err := sql.Open("postgres", l.connectionString.get())

	if err != nil {
		if l.debugPrint {
			log.Printf("couldn't connect database. %v\n", err)
		}

		return err
	}

	l.db = db

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(l.connectionTimeoutSeconds)*time.Second)

	defer cancel()

	if er2 := l.db.PingContext(ctx); er2 != nil {
		if l.debugPrint {
			log.Printf("couldn't ping database. %v\n", er2)
		}

		return er2
	}

	l.supposedReady = true

	return nil
}

func (l *DBLink) connectLoop() error {
	// If the first attempt failed, a more sophisticated approach is executed, addressing a controlled flow of sequential attempts
	var (
		connectionAttempt    uint
		connectionTimerStart = time.Now()
	)

	waitNextAttempt := func() {
		if l.debugPrint {
			log.Printf("wating %d seconds before retry to connect database.", l.timeBetweenConnectionAttemptsSeconds)
		}

		time.Sleep(time.Duration(l.timeBetweenConnectionAttemptsSeconds) * time.Second)
	}

	for {
		err := l.connectSimple()

		if err == nil {
			return nil
		}

		if time.Now().Before(connectionTimerStart.Add(time.Duration(l.connectionAttemptsMaxMinutes) * time.Minute)) {
			if l.debugPrint {
				log.Printf("database connection time limit reached without success. %v\n", err)
			}

			return fmt.Errorf("database connection time limit reached without success. %v", err)
		}

		connectionAttempt++

		if connectionAttempt >= l.connectionAttemptsMax {
			if l.debugPrint {
				log.Printf("database connection attempts limit reached without success. %v\n", err)
			}

			return fmt.Errorf("database connection attempts limit reached without success. %v", err)
		}

		waitNextAttempt()

		continue
	}
}

// Connect assures that database could be reached and used.
func (l *DBLink) Connect() error {
	if l.debugPrint {
		log.Println("opening database connection")
	}

	l.DBLinkOptions.print()

	if l.connectionAttemptsMax == 1 {
		return l.connectSimple()
	}

	return l.connectLoop()
}

// NewConnected returns an already connected instance of dbLink, ready for use
func NewConnected(opts *DBLinkOptions) (*DBLink, error) {
	dbl, err := New(opts)

	if err != nil {
		return nil, err
	}

	if er0 := dbl.Connect(); er0 != nil {
		return nil, er0
	}

	return dbl, nil
}

// Disconnect closes the link with database or database pool
func (l *DBLink) Disconnect() {
	if l.debugPrint {
		log.Println("disconnecting from database")
	}

	if l == nil || l.db == nil || !l.supposedReady {
		return
	}

	if err := l.db.Close(); err != nil {
		if l.debugPrint {
			log.Println(err)
		}
	}
}

// Database returns the raw sql driver database connection, aimed to let the programmer execute transactions and advanced operations
func (l *DBLink) Database() *sql.DB {
	return l.db
}

// IsEmptyErr returns true if the given error is sql.ErrNoRows
func (l *DBLink) IsEmptyErr(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, sql.ErrNoRows)
}

// New returns an instance of database, configured with options.
func New(opts *DBLinkOptions) (*DBLink, error) {
	if opts == nil {
		return nil, fmt.Errorf("undefined options")
	}

	if erv := opts.validate(); erv != nil {
		return nil, erv
	}

	dbl := DBLink{
		db:            nil,
		DBLinkOptions: *opts,
	}

	return &dbl, nil
}
