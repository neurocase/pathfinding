package main

import (
	"fmt"
	gameloop "github.com/GlenKelley/go-glutil/gameloop"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
		"github.com/neurocase/pathfinding"
		"math"
	//"time"
	//"math"
)


// define start, end and walls

//start, moves to square with lowest Hval that has not been visited

//		* give && check a list of available tiles
//		* keep a list of visited tiles (in order)

// there is the grid itself, then the nodes attached to the grid, although it has a position on the grid
// it is more specific to the path than the grid


/* each member of the grid array requires:
*  X,Y position value
*  TileType: (Wall, Open, Start Point, End Point)
*/


// each node requires :
/*	A Heuristic Value (manhattan based)
*	Distance moved (and order) from start
*	X, Y position
*/ 

// grid.movdist++ == ordered list
// !! ordered list - visited 
// if cannot move to new square, goto previous square and mark this square as a deadend 
// (repeat) untill found open unvisited tile.
// move to value with lowest hval 
// if at endpoint, return to start and goto path location with highest moved-distance

//repeat algorythm from start to points near to end path and calculate move distance
//see if heuristics from start path, to positions along path can be used to result in a lower move distance


/* In order to avoid confusion, avoid using 2D arrays, instead, use a 1d array
*  and subtract/add by line length in order to go to above/below grid array values
*/


var winHeight = 480.0
var winWidth = 480.0

var LineLength = 20

var mode = 0
var vx = 0.0
var vy = 0.0
var mousegridx = 0.0
var mousegridy = 0.0
var mousex = 0.0
var mousey = 0.0

var MazeMan = 21;
var StartPos = 21;
var EndPos = 389;

var grd = pathfinding.Entity{0,-10,0,0.1,"grey", true}
var here = 0

var Grid[400] int
var BreadCrumbs[]int
//var possiblelocations = []int{65,26,3,49,65,47,23,21}
var Inspection = []int {1,1,1,1,1,1,1,1}

var deadend = 0
var goback = 0


/*  0. open - grey
*	1. traveled - blue
*	2. wall - orange
*	3. deadend - red
*	4. startpos - green
*	5. endpos - purple 
*	6. thispos Lblue
*	7. mazeman White
*/	




func FindLowestCell(a []int)(int){
	deadend = 0
	lowest := 0
		for i := 0; i < 8; i++{

			if a[i] < a[lowest] {
				lowest = i
			}
		}
	
	if a[lowest] == 999{
	deadend = 1
	}
	return lowest

}



type XYpos struct{
	Xpos, Ypos int

}


type Node struct{
Gridnum, Hval int
IsStart bool
HaveTraveled bool

}
var Nodes[400] Node

func Inspect(x,y int)(int){
gnum :=	XYtoGrid(x, y)
rh := 999
if Grid[gnum] == 0 || Grid[gnum] == 1{
	rh = GiveHeuristic(gnum)
	Grid[gnum] = 1;
}
	return rh
}

func LookAround(pos int){
	Nodes[pos].HaveTraveled = true
	myx, myy := GridtoXY(pos)
	//Inspection := [8]int {}

	//up 123
	Inspection[0] = Inspect(myx-1,myy+1)
	Inspection[1] = Inspect(myx,myy+1)
	Inspection[2] = Inspect(myx+1,myy+1)

	//left right
	Inspection[3] = Inspect(myx-1,myy)
	Inspection[4] = Inspect(myx+1,myy)
	//down 123
	Inspection[5] = Inspect(myx-1,myy-1)
	Inspection[6] = Inspect(myx,myy-1)
	Inspection[7] = Inspect(myx+1,myy-1)
	
}

func GridAtDirection(m,d int)(int){
/*	0 up-left, 1 up, 2 up-right 
*	3 left, 4 right
*	6 down-left, 7 down, 8 down-right
*/
	myx, myy := GridtoXY(m)

	dest := 0

	switch d{
	case 0:
		dest = XYtoGrid(myx-1,myy+1)
	case 1:
		dest = XYtoGrid(myx,myy+1)
	case 2:
		dest = XYtoGrid(myx+1,myy+1)
	case 3:
		dest = XYtoGrid(myx-1,myy)
	case 4:
		dest = XYtoGrid(myx+1,myy)
	case 5:
		dest = XYtoGrid(myx-1,myy-1)
	case 6:
		dest = XYtoGrid(myx,myy-1)
	case 7:
		dest = XYtoGrid(myx+1,myy-1)
	}

	return dest
}




func GridtoXY(grd int) (int, int){

// 0 -> 19 = 0, | 20 -> 39 = 1

y := grd / LineLength -10
x := grd % LineLength -10

return x, y
}

func XYtoGrid(x, y int)(int){
x += 10
y += 10
grd := y*20 + x
return grd
}

func ResolutionToGrid(reswidth, resheight float64)(float64, float64){
ex := (reswidth * 2 / winWidth - 1)*10
wy := (1- resheight * 2 / winHeight)*10
return ex, wy
}

func WidthResToGrid(reswidth float64)(float64){
ex := (reswidth * 2 / winWidth - 1)*10
return ex
}

func HeightResToGrid(resheight float64)(float64){
wy := (1- resheight * 2 / winHeight)*10
//wy = -wy
return wy
}

func SetStartPos(st int){

	Grid[StartPos] = 0
	Grid[st] = 4
	StartPos = st
}

func SetManPos(man int){

	Grid[MazeMan] = 0
	Grid[man] = 6
	MazeMan = man
}

func SetEndPos(end int){

	Grid[EndPos] = 0
	Grid[end] = 5
	EndPos = end
}

func SwitchMode(){		
	if mode < 2{
		mode++
	}else{
		mode = 0
	}
}

func GiveHeuristic(grd int)(int){

	gx, gy := GridtoXY(grd)
	ex, ey := GridtoXY(EndPos)

	hx :=math.Abs(float64(gx) - float64(ex))
	hy :=math.Abs(float64(gy) - float64(ey))
	heur := hx+hy

	return int(heur)
}

func IsCrumbed(g int)(bool){

	for j := 0; j < len(BreadCrumbs)-1; j++{
			if g == BreadCrumbs[j]{
				return true
			}
	}
	return false
}


func main() {
	game := &Game{}
	err := gameloop.CreateWindow(int(winWidth), int(winHeight), "pathfinding gl test", false, game)
	fmt.Println(err)
}

type Game struct {
//	Red float64
}

func (game *Game) Init(window *glfw.Window) {
	//Select the 'projection matrix'
	gl.MatrixMode(gl.PROJECTION)
	//Reset
	gl.LoadIdentity()
	//Scale everything down, to 1/10 scale
	gl.Scaled(0.1,0.1,0.1)
	a,b := ResolutionToGrid(winWidth,winHeight)
	wrtg := WidthResToGrid(winWidth)
	hrtg := HeightResToGrid(winHeight)


	fmt.Println("Some Arbitary tests")

	fmt.Println("res to grid",winWidth, "=",a,"*", winHeight,"=",b)
	fmt.Println("Wres to grid", wrtg)
	fmt.Println("Hres to grid", hrtg)
	fmt.Println("Width to grid (100)", WidthResToGrid(100))
	fmt.Println("Height to grid (100)", HeightResToGrid(100))
	a,b = ResolutionToGrid(100,100)
	fmt.Println("Res to grid(100,100)",a,b)
	fmt.Println("")
	fmt.Println(">>>>>>> PRESS H FOR HELP <<<<<<<")
}

func (game *Game) Draw(window *glfw.Window) {
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.ClearColor(0.2, 0.2, 0.2, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	//mousegridx := 0.0
	//mousegridy := 0.0
	mousegridx, mousegridy = ResolutionToGrid(mousex,mousey)
	here = XYtoGrid(int(mousegridx), int(mousegridy))
	
	for i := 0 ; i < 400; i++{
				
		if IsCrumbed(i){
			Grid[i] = 7
		}
		if i == MazeMan{
			Grid[i] = 6
		}
		

		switch Grid[i]{
				case 0:
					grd.Colour = "grey"
					grd.Size = 0.1
				case 1:
					grd.Colour = "blue"
					grd.Size = 0.3
				case 2:
					grd.Colour = "orange"
					grd.Size = 0.3
				case 3:
					grd.Colour = "red"
					grd.Size = 0.2
				case 4:
					grd.Colour = "green"
					grd.Size = 0.5
				case 5:
					grd.Colour = "purple"
					grd.Size = 0.5
				case 6:
					grd.Colour = "white"
					grd.Size = 0.5
				case 7:
					grd.Colour = "lblue"
					grd.Size = 0.3
		}	

		gx, gy := GridtoXY(i)
		grd.Xpos = float64(gx)
		grd.Ypos = float64(gy)
		pathfinding.DrawEntity(grd)



	}
	//func XYtoGrid(x, y int)(int){
	//func GridtoXY(grd int) (int, int){

	mousegridx, mousegridy = ResolutionToGrid(mousex,mousey)
	redtri := pathfinding.Entity{0,0,0,0.5,"yellow", true}
	redtri.Xpos = float64(int(mousegridx))
	redtri.Ypos = float64(int(mousegridy))
			pathfinding.DrawEntity(redtri)
	
	//gx := int(vx)+10
	//gy := int(vy)+10

}

func (game *Game) Reshape(window *glfw.Window, width, height int) {
}

func (game *Game) MouseClick(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	if action == glfw.Press{

		switch button{

		case glfw.MouseButtonLeft:

			fmt.Println(mousegridx,mousegridy)

			//click := XYtoGrid(int(mousegridx), int(mousegridy))

			fmt.Println("painting tile number:", here, "at",mousegridx,mousegridy)

			switch mode{
				case 0:
					Grid[here] = 2
				case 1:
					SetStartPos(here)
					if MazeMan == 0{
						SetManPos(here)
					}
				case 2:
					SetEndPos(here)
			}
			
		case glfw.MouseButtonRight:

			fmt.Println(mousegridx,mousegridy)

			fmt.Println(here)
			Grid[here] = 0

		}
	}
}



func (game *Game) MouseMove(window *glfw.Window, xpos float64, ypos float64) {
	mousex = xpos
	mousey = ypos
	//pos := XYtoGrid(int(mousegridx), int(mousegridy))
}

func (game *Game) KeyPress(window *glfw.Window, k glfw.Key, s int, action glfw.Action, mods glfw.ModifierKey) {
	//fmt.Println("keypress", k)
	if action == glfw.Release {
		switch k {

		case glfw.KeyEscape:
			window.SetShouldClose(true)
			//b key
		case  66:
			fmt.Println((len(BreadCrumbs)-1), "Crumz")
			//if len(BreadCrumbs) < 1{
				for i := 0; i < len(BreadCrumbs)-1;i++{
					fmt.Println(i,":",BreadCrumbs[i])
			//	}
			}
			// > key
		case 46:
			LookAround(MazeMan)
			inspdir := FindLowestCell(Inspection)
			if deadend == 0{
			dir := GridAtDirection(MazeMan,inspdir)
			SetManPos(dir)
			BreadCrumbs = append(BreadCrumbs,dir)
			goback=0
			}else{
			fmt.Println("deadend, go back")
			b := len(BreadCrumbs)
			c := BreadCrumbs[b-1-goback]
			goback++
			SetManPos(c)
			}
			
		case 72:
			fmt.Println("Z: RESET|X: HEURISTIC|C: COMPUTEPATH|M:BRUSH|V:NEXTPOINT")

			//X key
		case 88:
				heur := GiveHeuristic(here)
				fmt.Println("Here:",here,"Heuristic:",heur, "Crumbs:", IsCrumbed(here))


			//v Key
		case 86:
				inspdir := FindLowestCell(Inspection)
				dir := GridAtDirection(MazeMan,inspdir)
				fmt.Println("next cell is",dir)
			//Z KEY
		case 90:
			BreadCrumbs = BreadCrumbs[:0]
			for i := 0; i < 400; i++{
				if Grid[i] == 1{
					Grid[i] = 0
				}
			}
			SetManPos(StartPos)
			fmt.Println("MazeMan sent to StartPos")

			//C KEY
		case 67:
			LookAround(MazeMan)
			fmt.Println("ComputePath")

				sortloc := FindLowestCell(Inspection)
				fmt.Println(sortloc)

			//M KEY
		case 77:
			SwitchMode()
			modeis := "wall"
			switch mode{
				case 0:
					modeis = "wall"
				case 1:
					modeis = "start position"
				case 2: 
					modeis = "end position"
			}
			fmt.Println(modeis, "brush")
		//case 46:
			//fmt.Println("MOVE FORWARDS")
		case 44:
			fmt.Println("MOVE BACKWARDS")
		}
	}
}

func (game *Game) Scroll(window *glfw.Window, xoff float64, yoff float64) {

}

func (game *Game) Simulate(time gameloop.GameTime) {
	//game.Red = math.Sin(time.Elapsed.Seconds())
}

func (game *Game) OnClose(window *glfw.Window) {

}

func (game *Game) IsIdle() bool {
	//if idle is true, the gameloop will
	//wait for user input before drawing
	return false
}

func (game *Game) NeedsRender() bool {
	//if render is false the game will not redraw the screen
	return true
}
