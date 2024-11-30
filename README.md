# code-words
Learning Go by creating some games to run on a terminal


## Todo

1. Vectors and Rectangles
2. Logger + Debug line
3. Raw mode terminal
4. Input handling and exiting

First, get collision logic working with just points and fields
    Print to console, don't worry about raw mode
    Have dots move along a pre-programmed path
Second, add a logger to record inputs + time offsets
Third, way to define and load up a game's history and play using those input timings + apple spawns

Buffer work as just fields of bytes that are overlaid on top of the layer below
Special case things w/ transparency enabled will not overwrite layers below if they're empty
    In this case, use an abstracted grid above the byte layer
    Cells in the grid can have one of N types
    Use grid to determine collisions
    Linkedlist of cell indices represents the snake
    Render function will iterate over that grid to assemble the byte array

## Future Plans
new plan for the game. Snake text editor
"Apples" define a byte or rune that is used for display
    Could also use ghostly "'" characters to delineate blank space, newline, etc.
    Maybe even have multi-character elements that you can eat in one go
Screen split into two sections with a door down the middle

---------------------------------------------------------------------------------------
|This is the target text (written in gray)      |    F                                |
|This is the actual text (written in white)     |               E     5               |
|                                               |                                     |
|<targetting this line to add>                           6                            |
|                                     +ABCDEF0123456789ABC                            |
|                                                        D                            |
|                                               |        EF0     4                    |
|                                               |        1                            |
|                                               |               D           7         |
|                                               |                                     |
|                                               |    3                                |
|                                               |                                     |
|                                               |        C       2        8           |
|                                               |                                     |
|                                                                                     |
|                                                          B                          |
|                                                                                     |
|                                               |           A    0    9               |
|                                               |                                     |
|                                               |                                     |
=======================================================================================
|+this is the actual text form of your input bytes                                    |
=======================================================================================


Get points for retrieving characters
Scrabble bonuses for completing certain words (Foreground/background colours?)
Word + character multipliers for adding a certain section
    - Multiplications if you proc multiple
    - Word length is base
    - other additions/multiplications apply left to right
    - additional distance in multiplier, adding the last bit is very good

Maybe switch to something w/ rasterized graphics to make it easier to display complex behaviour
Add support for requesting nearest apple of type to point



## System Design

This is meant to be a lobby based experience. What sort of modes should be available?
It'd be cool to set up a P2P mode for just terminals and a separate server based mode for doing a web based model so I could get used to HTML, HTTP, and general web animations

Step one, simple TCP connection between two terminals
    // Was able to set up a three connections with multiple messages being sent and received between them
    // next part here is going to be formalizing the server/client setup so I don't need to do it myself every time
    // Then I'll need to break down the client logic vs the server logic into something that can be kicked off independently
    // Last thing is going to be a convenience thing that will fire off a server and multiple clients and gather the logs from them all.
        // Maybe it'll just work and append nicely. We'll see...
Step two, One to many connection with a single computer serving as the "server" / main node, then other clients that are able to connect and say hello
Step three, set up flows for initiating a new connection and validating messages received (req/ack, rep/ack)
    // May well encounter issues, but will roll with this for at least v1
Step four, define requirements for messaging protocol based off of this
    // Want to minimize overhead + size of messages, the rest I don't know
Step five, testing logic for sending and receiving bulk messages + handling multiple messages concurrently
    // Don't need to process concurrently, just need to handle it gracefully
Step six, first pass logic at actually sending and receiving state
    // What sort of metrics do I care about?
    // This is all on local network, so not worried about sending a lot of packets
    // Main concern is ensuring that we are able to communicate quickly and parse things out in a minimal amount of time
Step seven, Client side rendering, given a state packet, parse it out and colour it in