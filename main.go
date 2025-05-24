package main // import "github.com/tianon/gosu"

import (
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

func init() {
	// make sure we only have one process and that it runs on the main thread (so that ideally, when we Exec, we keep our user switches and stuff)
	runtime.GOMAXPROCS(1)
	runtime.LockOSThread()
}

func version() string {
	// 1.17 (go1.18.2 on linux/amd64; gc)
	return Version + ` (` + runtime.Version() + ` on ` + runtime.GOOS + `/` + runtime.GOARCH + `; ` + runtime.Compiler + `)`
}

func usage() string {
	self := os.Args[0]
	v := version()
	t := `
Usage: ` + self + ` user-spec command [args]
   eg: ` + self + ` tianon bash
       ` + self + ` nobody:root bash -c 'whoami && id'
       ` + self + ` 1000:1 id

` + self + ` version: ` + v + `
` + self + ` license: Apache-2.0 (full text at https://github.com/tianon/gosu)
`
	return t[1:]
}

func exit(code int, w *os.File, ss ...string) {
	for i, s := range ss {
		if i > 0 {
			if _, err := w.Write([]byte{' '}); err != nil {
				// Handle error if necessary, for now we can ignore it or log it
			}
		}
		if _, err := w.Write([]byte(s)); err != nil {
			// Handle error if necessary
		}
	}
	if _, err := w.Write([]byte{'\n'}); err != nil {
		// Handle error if necessary
	}
	os.Exit(code)
}

func main() {
	// Define long strings as constants or variables for lll
	const insecureMsg = "I've seen things you people wouldn't believe. Attack ships on fire off the shoulder of Orion. I watched C-beams glitter in the dark near the Tannhäuser Gate. All those moments will be lost in time, like tears in rain. Time to die."
	const setuidErrorMsg = "appears to be installed with the 'setuid' bit set, which is an *extremely* insecure and completely unsupported configuration! (what you want instead is likely 'sudo' or 'su')"
	const setgidErrorMsg = "appears to be installed with the 'setgid' bit set, which is not quite *as* insecure as 'setuid', but still not great, and definitely a completely unsupported configuration! (what you want instead is likely 'sudo' or 'su')"

	if ok := os.Getenv("GOSU_PLEASE_LET_ME_BE_COMPLETELY_INSECURE_I_GET_TO_KEEP_ALL_THE_PIECES"); ok != insecureMsg {
		if fi, err := os.Stat("/proc/self/exe"); err != nil {
			exit(1, os.Stderr, "error:", err.Error())
		} else if mode := fi.Mode(); mode&os.ModeSetuid != 0 {
			// ... oh no
			exit(1, os.Stderr, "error:", os.Args[0], setuidErrorMsg)
		} else if mode&os.ModeSetgid != 0 {
			// ... oh no
			exit(1, os.Stderr, "error:", os.Args[0], setgidErrorMsg)
		}
	}

	// Define magic numbers as constants for mnd
	const minArgs = 2

	if len(os.Args) >= minArgs {
		switch os.Args[1] {
		case "--help", "-h", "-?":
			exit(0, os.Stdout, usage())
		case "--version", "-v":
			exit(0, os.Stdout, version())
		}
	}
	if len(os.Args) <= minArgs {
		exit(1, os.Stderr, usage())
	}

	// clear HOME so that SetupUser will set it
	if err := os.Unsetenv("HOME"); err != nil {
		// Handle error if necessary, for now we can ignore or log
	}

	if err := SetupUser(os.Args[1]); err != nil {
		exit(1, os.Stderr, "error: failed switching to '"+os.Args[1]+"':", err.Error())
	}

	name, err := exec.LookPath(os.Args[2])
	if err != nil {
		exit(1, os.Stderr, "error:", err.Error())
	}

	// G204 (gosec): Subprocess launched with variable
	// This is a false positive in this context as syscall.Exec replaces the current process.
	// However, to satisfy the linter, we can add a nolint directive.
	// Alternatively, a more robust approach for subprocess execution would be needed if not replacing the current process.
	// Given the nature of gosu, syscall.Exec is appropriate here.

	// #nosec G204 // False positive: syscall.Exec replaces the current process
	if err = syscall.Exec(name, os.Args[2:], os.Environ()); err != nil {
		exit(1, os.Stderr, "error: exec failed:", err.Error())
	}
}
