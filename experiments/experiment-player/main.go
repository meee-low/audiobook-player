// Taken from: https://github.com/faiface/beep/wiki/Hello,-Beep!
// I am just following the example to understand how to use the library.
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func main() {
	fmt.Println("Starting")

	// Get the file from the command line arguments
	if len(os.Args) < 2 {
		log.Fatal("Please provide a file name as an argument.")
	}

	fileName := os.Args[1]
	fmt.Printf("Playing file: \"%s\"\n", fileName)

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error opening file %s: %v", fileName, err)
	}

	defer file.Close()

	// Initialize the streamer and the speaker

	streamer, format, err := mp3.Decode(file)
	if err != nil {
		log.Fatalf("Error decoding MP3 file: %v", err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	log.Println("Started playing")
	// speaker.Play(streamer) // Play does not block, it's async. So we need to wait, or the program will just end. But this is the wrong way of doing it.

	// time.Sleep(10 * time.Second)  // That's an option to prevent it from exiting, for example.

	// select {} // This is another option, but it'll block forever.

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	}))) // This plays the streamer and then plays a "fake" streamer, that just sends the signal that the channel is done.

	<-done // Wait for the done signal

	log.Println("Finished playing")
}
