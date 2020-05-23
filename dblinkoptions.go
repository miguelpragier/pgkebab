package pgkebab

import "fmt"

// DBLinkOptions contains directives for config database connection
// To create an instance, use Options() function
type DBLinkOptions struct {
	// connectionAttemptsMax determines how many attempts the engine should do before give away with an error.
	// Zero means retry indefinitely or until connectionAttemptsMaxMinutes is reached
	connectionAttemptsMax uint
	// connectionAttemptsMaxMinutes determines how much time the engine keeps trying to (re)connect before fail with an error
	// Zero means retry indefinitely or until connectionAttemptsMax is reached
	connectionAttemptsMaxMinutes uint
	// timeBetweenConnectionAttemptsSeconds - How much time between connection attempts
	timeBetweenConnectionAttemptsSeconds uint
	// database connectionString
	connectionString *ConnectionString
	// connectionTimeoutSeconds Seconds to wait for (re)connection
	connectionTimeoutSeconds uint
	// executionTimeoutSeconds seconds for wait queries execution
	executionTimeoutSeconds uint
	// maxOpenConnections Max Open Connections
	//	maxOpenConnections int
	// emergencyCallback A routine to call when the database can't be reached
	emergencyCallback func(error)
	// debugPrint defines if detailed information about errors and general execution gonna be printed
	debugPrint bool
}

func (o *DBLinkOptions) validate() error {
	if o.connectionString == nil {
		return fmt.Errorf("empty connection string")
	}

	if err := o.connectionString.validate(); err != nil {
		return err
	}

	if o.timeBetweenConnectionAttemptsSeconds < 1 {
		o.timeBetweenConnectionAttemptsSeconds = timeBetweenConnectionAttemptsSecondsDefault
	}

	if o.connectionTimeoutSeconds == 0 {
		o.connectionTimeoutSeconds = connectionTimeoutSecondsDefault
	}

	if o.executionTimeoutSeconds == 0 {
		o.executionTimeoutSeconds = executionTimeoutSecondsDefault
	}

	return nil
}

func (o DBLinkOptions) print() {

	if !o.debugPrint {
		return
	}

	fmt.Println("-- ---- ---- ---- ---- ---- ---- ---- ---- ---- ---- ---- ---- ---- ---- ---- ---- ---- --")
	fmt.Println(`PGKebab Options:`)
	fmt.Printf("-- Max (Re)Connection Attempts: %d\n", o.connectionAttemptsMax)
	fmt.Printf("-- Max (Re)Connection Attempts in minutes: %d\n", o.connectionAttemptsMaxMinutes)
	fmt.Printf("-- Seconds between attempts: %d\n", o.timeBetweenConnectionAttemptsSeconds)
	fmt.Printf("-- Connection Timeout in seconds: %d\n", o.connectionTimeoutSeconds)
	fmt.Printf("-- Execution timeout in seconds: %d\n", o.executionTimeoutSeconds)
	fmt.Println("-- If you can see this summary, pgKebab was initialized with [debugPrint] option on")
	fmt.Println("-- ---- ---- ---- ---- ---- ---- ---- ---- ---- ---- ---- ---- ---- ---- ---- ---- ---- --")
}

// Options is a helper to initializate a new DBLinkOptions structure
// reconnectionAttemptsMax is the max number of (re)connection attempts. If zeroed, keeps trying undefinitely or until reconnectionAttemptsMaxMinutes is reached.
// reconnectionAttemptsMaxMinutes is the max allowed time to insist in (re)connect. If zeroed, keeps trying undefinetely or until reconnectionAttemptsMax is reached
// intervalBetweenReconnectionAttemptsSeconds is the sleepin' time lapse between (re)connection attempts
// debugPrint if true, prints debug/log messages to stdout with stdlib log.Printf() function
func Options(cs *ConnectionString, connTimeout, execTimeout, connAttemptsMax, connAttemptsMaxMinutes, secondsBetweenReconnectionAttempts uint, debugPrint bool) *DBLinkOptions {
	return &DBLinkOptions{
		connectionString:                     cs,
		connectionTimeoutSeconds:             connTimeout,
		executionTimeoutSeconds:              execTimeout,
		connectionAttemptsMax:                connAttemptsMax,
		connectionAttemptsMaxMinutes:         connAttemptsMaxMinutes,
		timeBetweenConnectionAttemptsSeconds: secondsBetweenReconnectionAttempts,
		debugPrint:                           debugPrint,
		emergencyCallback:                    nil,
	}
}
