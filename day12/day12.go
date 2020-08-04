package main

import (
	"fmt"
)

type positionType struct {
	X int
	Y int
	Z int
}

type moonType struct {
	pos positionType
	vel positionType
}

type moonSet []moonType

var (
	debug      = false
	itteration = 0
	done       = make(chan int)
)

func main() {

	moons := moonSet{

		// example  set
		// {pos: positionType{X:  -1, Y:   0, Z:  2}, vel: positionType{0, 0, 0}},
		// {pos: positionType{X:   2, Y: -10, Z: -7}, vel: positionType{0, 0, 0}},
		// {pos: positionType{X:   4, Y:  -8, Z:  8}, vel: positionType{0, 0, 0}},
		// {pos: positionType{X:   3, Y:   5, Z: -1}, vel: positionType{0, 0, 0}},

		// example set two 
		// {pos: positionType{X:  -8, Y: -10, Z:   0}, vel: positionType{0, 0, 0}},
		// {pos: positionType{X:   5, Y:   5, Z:  10}, vel: positionType{0, 0, 0}},
		// {pos: positionType{X:   2, Y:  -7, Z:   3}, vel: positionType{0, 0, 0}},
		// {pos: positionType{X:   9, Y:  -8, Z:  -3}, vel: positionType{0, 0, 0}},
 
		// problem set
		{ pos: positionType{X:  9, Y: 13, Z:  -8}, vel: positionType{0,0,0} },
		{ pos: positionType{X: -3, Y: 16, Z: -17}, vel: positionType{0,0,0} },
		{ pos: positionType{X: -4, Y: 11, Z: -10}, vel: positionType{0,0,0} },
		{ pos: positionType{X:  0, Y: -2, Z:  -2}, vel: positionType{0,0,0} },
	}

	// moons.print()
	// fmt.Println()

	for i := 0; i < 1000; i++ {
		// fmt.Printf("\nItteration:%d\n", i)
		moons.adjustVelocity()
		moons.adjustPosition()
		// moons.print()
	}

	sum := 0
	for i := 0; i < len(moons); i++ {
		sum += (moons[i]).energy()
	}
	fmt.Printf("Total Energy: %d",sum)
}

func (moons moonSet) adjustVelocity() {
	for i := 0; i < len(moons)-1; i++ {
		for j := i + 1; j < len(moons); j++ {
			// fmt.Printf("i: %d j: %d  -- moonI:%v moonJ:%v \n", i, j, moons[i], moons[j])
			if moons[i].pos.X < moons[j].pos.X {
				moons[i].vel.X++
				moons[j].vel.X--
			}
			if moons[i].pos.X > moons[j].pos.X {
				moons[i].vel.X--
				moons[j].vel.X++
			}

			if moons[i].pos.Y < moons[j].pos.Y {
				moons[i].vel.Y++
				moons[j].vel.Y--
			}
			if moons[i].pos.Y > moons[j].pos.Y {
				moons[i].vel.Y--
				moons[j].vel.Y++
			}

			if moons[i].pos.Z < moons[j].pos.Z {
				moons[i].vel.Z++
				moons[j].vel.Z--
			}
			if moons[i].pos.Z > moons[j].pos.Z {
				moons[i].vel.Z--
				moons[j].vel.Z++
			}
			// fmt.Println()
		}
	}
}

func (moons moonSet) adjustPosition() {
	for i := 0; i < len(moons); i++ {
		moons[i].pos.X += moons[i].vel.X
		moons[i].pos.Y += moons[i].vel.Y
		moons[i].pos.Z += moons[i].vel.Z
	}
}

func (moons moonSet) print() {
	for i := 0; i < len(moons); i++ {
		fmt.Printf("%v\n", moons[i])
	}
}

func (m moonType) energy() int {
	pe := (abs(m.pos.X) + abs(m.pos.Y) + abs(m.pos.Z))
	ke := (abs(m.vel.X) + abs(m.vel.Y) + abs(m.vel.Z))
	return pe * ke
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
