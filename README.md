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