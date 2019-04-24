package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

var (
	errInputError = fmt.Errorf("Function input error")
)

// Point represents point
type Point []float64

// Clusters represents set of clusters
type Clusters map[int][]Point

// NewPoint creates new point
func NewPoint(args ...float64) (result Point) {
	result = make(Point, len(args))
	for i, elem := range args {
		result[i] = elem
	}

	return
}

// getDistance returns the distance between points
func getDistance(a Point, b Point) (result float64) {
	for i := range a {
		result += math.Pow(a[i]-b[i], 2)
	}

	return math.Pow(result, 0.5)
}

// getHypersphere возвращает все точки из массива, которые находятся в гиперсфере
// указанного радиуса из указанным центром
//
// points - Массив точек
// center - Центр гиперсферы
// radius - Радиус гиперсферы
//
func getHypersphere(points []Point, center Point, radius float64) (result []Point) {
	result = make([]Point, 0)
	indexes := make([]int, 0)

	for i, point := range points {
		if getDistance(point, center) < radius {
			result = append(result, point)
			indexes = append(indexes, i)
		}
	}

	return
}

//
//	findCenter находит центр масс набора точек
//
//	points - набор точек
//
func findCenter(points []Point) (result Point) {
	dimensions := len(points[0])
	result = make(Point, dimensions)

	for i := 0; i < dimensions; i++ {
		for _, point := range points {
			result[i] += point[i]
		}

		result[i] /= float64(len(points))
	}

	return
}

func getMinBound(points []Point) (result Point) {
	dimensions := len(points[0])
	result = make(Point, dimensions)
	copy(result, points[0])

	for _, points := range points {
		for i := 0; i < dimensions; i++ {
			if points[i] < result[i] {
				result[i] = points[i]
			}
		}
	}

	return
}

func getMaxBound(points []Point) (result Point) {
	dimensions := len(points[0])
	result = make(Point, dimensions)
	copy(result, points[0])

	for _, points := range points {
		for i := 0; i < dimensions; i++ {
			if points[i] > result[i] {
				result[i] = points[i]
			}
		}
	}

	return
}

func normPoint(point Point, minBound Point, maxBound Point) (result Point) {
	dimensions := len(point)
	result = make(Point, dimensions)

	for i := 0; i < dimensions; i++ {
		result[i] = (point[i] - minBound[i]) / (maxBound[i] - minBound[i])
	}

	return
}

func denormPoint(normPoint Point, minBound Point, maxBound Point) (result Point) {
	dimensions := len(normPoint)
	result = make(Point, dimensions)

	for i := 0; i < dimensions; i++ {
		result[i] = normPoint[i]*(maxBound[i]-minBound[i]) + minBound[i]
	}

	return
}

func comparePoints(a Point, b Point) bool {
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func deletePoint(points []Point, point Point) (result []Point) {
	for i, elem := range points {
		if comparePoints(elem, point) {
			result = append(points[:i], points[i+1:]...)
			return
		}
	}

	return points
}

func validatePoints(points []Point) (result bool) {
	if len(points) == 0 {
		return false
	}

	dimensions := len(points[0])
	if dimensions == 0 {
		return false
	}

	for _, point := range points {
		if len(point) != dimensions {
			return false
		}
	}

	return true
}

// Clusterize запускает кластеризацию точек
//
// points      - список точек
// minClusters - минимальное возможное количество кластеров,
//               если равно 0, то ограничения не будет
func Clusterize(points []Point, minClusters int) (clusters Clusters, err error) {
	if !validatePoints(points) {
		return nil, errInputError
	}

	minCoords := getMinBound(points)
	maxCoords := getMaxBound(points)

	normPoints := make([]Point, 0)
	normPointsSave := make([]Point, 0)

	// Нормирование
	for _, point := range points {
		normPoint := normPoint(point, minCoords, maxCoords)
		normPoints = append(normPoints, normPoint)

		// Так как точки в norm_points будут удаляться, сюда сохранено исходное состояние
		normPointsSave = append(normPointsSave, normPoint)
	}

	// начальное значение радиуса гиперсферы
	dimensions := len(points[0])
	radius := math.Pow(float64(dimensions), 0.5) / 2
	var k float64

	for {
		clusters = make(Clusters)
		radius -= radius * (k + 1) / 10
		i := 0

		for len(normPoints) > 0 {
			centers := make([]Point, 0)

			sphere := getHypersphere(normPoints, normPoints[0], radius)
			centers = append(centers, findCenter(sphere))

			sphere = getHypersphere(normPoints, centers[0], radius)
			centers = append(centers, findCenter(sphere))

			for getDistance(centers[len(centers)-1], centers[len(centers)-2]) > 0.0005 {
				sphere = getHypersphere(normPoints, centers[len(centers)-1], radius)
				centers = append(centers, findCenter(sphere))
			}

			cluster := make([]Point, 0)
			for _, point := range sphere {
				cluster = append(cluster, denormPoint(point, minCoords, maxCoords))
				normPoints = deletePoint(normPoints, point)
			}

			clusters[i] = cluster
			i++
		}
		k++

		if minClusters <= 0 {
			return clusters, nil
		}

		if len(clusters) >= minClusters {
			return clusters, nil
		}

		normPoints = make([]Point, 0)
		for _, elem := range normPointsSave {
			normPoints = append(normPoints, elem)
		}
	}
}

// Чтение входных данных из файла
//
// fileName    - путь к входному файлу
// minClusters - по этому адресу запишется минимальное количество кластеров
func readFromFile(fileName string, minClusters *int) (result []Point, err error) {
	file, err := os.Open(fileName)
	start := true
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	result = make([]Point, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if start {
			var count int
			count, err = strconv.Atoi(scanner.Text())
			if err != nil {
				result = nil
				return
			}

			*minClusters = count
			start = false
			continue
		}

		coordsStr := strings.Split(scanner.Text(), ",")
		coords := make(Point, 0)

		for _, coordStr := range coordsStr {
			var coord float64
			coord, err = strconv.ParseFloat(coordStr, 64)
			if err != nil {
				result = nil
				return
			}
			coords = append(coords, coord)

		}

		result = append(result, coords)
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return
}

func main() {
	var minClusters int
	points, err := readFromFile("input.csv", &minClusters)
	if err != nil {
		fmt.Println(err)
		return
	}

	clusters, err := Clusterize(points, minClusters)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, cluster := range clusters {
		for _, point := range cluster {
			fmt.Printf("%v\n", point)
		}
		fmt.Println("-------")
	}
}
