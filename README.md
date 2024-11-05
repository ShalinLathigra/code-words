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