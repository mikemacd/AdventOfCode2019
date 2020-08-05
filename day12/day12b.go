package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

/*
           4686774924
Itteration: 670000000
Itteration:1745000000
           4686774924
Itteration:  4332000000
Itteration:134004000000
Itteration:134 004 000 000
Itteration:200178000000

*/

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

type moonCollection struct {
	m0 moonType
	m1 moonType
	m2 moonType
	m3 moonType
}

type moonIndex struct {
	posA int
	velA int
	posB int
	velB int
	posC int
	velC int
	posD int
	velD int
}

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
		{pos: positionType{X: 9, Y: 13, Z: -8}, vel: positionType{0, 0, 0}},
		{pos: positionType{X: -3, Y: 16, Z: -17}, vel: positionType{0, 0, 0}},
		{pos: positionType{X: -4, Y: 11, Z: -10}, vel: positionType{0, 0, 0}},
		{pos: positionType{X: 0, Y: -2, Z: -2}, vel: positionType{0, 0, 0}},
	}

	_=spew.Sdump("")
 
	fX, fY, fZ := 0, 0, 0
	fX = calcX(moons)
	fY = calcY(moons)
	fZ = calcZ(moons)
	fmt.Printf("%v %v %v\n", fX, fY, fZ)

	lcmXYZ := lcm(fX,fY,fZ)

	fmt.Printf("LCM:%d\n", lcmXYZ)
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
	// fmt.Println()
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

func (moons moonSet) equal(original moonSet) bool {
	equal := true
	for i := 0; i < len(moons); i++ {
		equal = equal && (moons[i] != original[i])
	}
	return equal
}

//func (m moonType)

func calcX(moons moonSet) int {
	seen := map[moonIndex]int{}

	i := -1

	for {
		i++

		for i := 0; i < len(moons)-1; i++ {
			for j := i + 1; j < len(moons); j++ {
				if moons[i].pos.X < moons[j].pos.X {
					moons[i].vel.X++
					moons[j].vel.X--
				}
				if moons[i].pos.X > moons[j].pos.X {
					moons[i].vel.X--
					moons[j].vel.X++
				}
			}
		}
		for i := 0; i < len(moons); i++ {
			moons[i].pos.X += moons[i].vel.X
		}

		mc := moonIndex{
			moons[0].pos.X, moons[0].vel.X, moons[1].pos.X, moons[1].vel.X,
			moons[2].pos.X, moons[2].vel.X, moons[3].pos.X, moons[3].vel.X,
		}
		if _, ok := seen[mc]; !ok {
			seen[mc] = 0
		}
		seen[mc]++

		// fmt.Printf("i: %d j: %d  -- moonI:%v moonJ:%v \n", i, j, moons[i], moons[j])

		// moons.print()
		if seen[mc] > 1 {
			return i
			fmt.Printf("Original position after %d itterations:  \n", i)
			break
		}
	}
	return -1
}
func calcY(moons moonSet) int {
	seen := map[moonIndex]int{}

	i := -1

	for {
		i++

		for i := 0; i < len(moons)-1; i++ {
			for j := i + 1; j < len(moons); j++ {
				if moons[i].pos.Y < moons[j].pos.Y {
					moons[i].vel.Y++
					moons[j].vel.Y--
				}
				if moons[i].pos.Y > moons[j].pos.Y {
					moons[i].vel.Y--
					moons[j].vel.Y++
				}
			}
		}
		for i := 0; i < len(moons); i++ {
			moons[i].pos.Y += moons[i].vel.Y
		}

		mc := moonIndex{
			moons[0].pos.Y, moons[0].vel.Y, moons[1].pos.Y, moons[1].vel.Y,
			moons[2].pos.Y, moons[2].vel.Y, moons[3].pos.Y, moons[3].vel.Y,
		}
		if _, ok := seen[mc]; !ok {
			seen[mc] = 0
		}
		seen[mc]++

		// fmt.Printf("i: %d j: %d  -- moonI:%v moonJ:%v \n", i, j, moons[i], moons[j])

		// moons.print()
		if seen[mc] > 1 {
			return i
			fmt.Printf("Original position after %d itterations:  \n", i)
			break
		}
	}
	return -1

}
func calcZ(moons moonSet) int {
	seen := map[moonIndex]int{}

	i := -1

	for {
		i++

		for i := 0; i < len(moons)-1; i++ {
			for j := i + 1; j < len(moons); j++ {
				if moons[i].pos.Z < moons[j].pos.Z {
					moons[i].vel.Z++
					moons[j].vel.Z--
				}
				if moons[i].pos.Z > moons[j].pos.Z {
					moons[i].vel.Z--
					moons[j].vel.Z++
				}
			}
		}
		for i := 0; i < len(moons); i++ {
			moons[i].pos.Z += moons[i].vel.Z
		}

		mc := moonIndex{
			moons[0].pos.Z, moons[0].vel.Z, moons[1].pos.Z, moons[1].vel.Z,
			moons[2].pos.Z, moons[2].vel.Z, moons[3].pos.Z, moons[3].vel.Z,
		}
		if _, ok := seen[mc]; !ok {
			seen[mc] = 0
		}
		seen[mc]++

		// fmt.Printf("i: %d j: %d  -- moonI:%v moonJ:%v \n", i, j, moons[i], moons[j])

		// moons.print()
		if seen[mc] > 1 {
			return i
			fmt.Printf("Original position after %d itterations:  \n", i)
			break
		}
	}
	return -1

}
func gcd(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func lcm(a, b int, integers ...int) int {
	result := a * b / gcd(a, b)

	for i := 0; i < len(integers); i++ {
		result = lcm(result, integers[i])
	}

	return result
}
