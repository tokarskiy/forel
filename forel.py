#################################################
### Реализация алгоритма кластеризации FOREL
#################################################
import csv

###
### Реализация алгоритма FOREL
###
### points       - список точек, точка - это вектор её координат
### min_clusters - минимальное возможное количество кластеров
###
def forel(points, min_clusters):
    dimensions = len(points[0])

    # Найти минимум и максимум по признакам 
    min_coords = [min_list([points[j][i] for j in range(0, len(points))]) for i in range(0, dimensions)]
    max_coords = [max_list([points[j][i] for j in range(0, len(points))]) for i in range(0, dimensions)]

    # Нормирование 
    norm_points = []
    norm_points_save = [] # Так как точки в norm_points будут удаляться, сюда сохранено исходное состояние
    for point in points:
        norm_point = [(point[i] - min_coords[i]) / (max_coords[i] - min_coords[i]) for i in range(0, dimensions)]
        norm_points.append(norm_point)
        norm_points_save.append(norm_point)
    

    radius = (dimensions ** 0.5) / 2 # начальное значение радиуса гиперсферы
    clusters = []
    k = 0
    while True:
        clusters = []
        radius -= radius * (k + 1) / 10
        while len(norm_points) > 0:
            centers = []

            sphere = get_hypersphere(norm_points, norm_points[0], radius)
            centers.append(center_list(sphere))

            sphere = get_hypersphere(norm_points, centers[0], radius)
            centers.append(center_list(sphere))
        
            # Пока разница между центрами на текущей и предыдущей итерации меньше чем эпсилон 
            while distance(centers[-1], centers[-2]) > 0.0005:
                sphere = get_hypersphere(norm_points, centers[-1], radius)
                centers.append(center_list(sphere))

            cluster = []
            # Запись результата в кластер
            for point in sphere:
                # Восстановление исходной точки из нормированной 
                new_point = [point[i] * (max_coords[i] - min_coords[i]) + min_coords[i] for i in range(dimensions)]
                cluster.append(new_point)

            # Запись кластера в массив кластеров и удаление "занятых" точек
            clusters.append(cluster)
            for point in sphere:
                norm_points.remove(point)
        
        k = k + 1
        for point in norm_points_save:
            norm_points.append(point)

        # В случае, если количество кластеров меньше, чем min_clusters, проходим алгоритм заново с
        # уменьшенным радиусом гиперсферы
        if len(clusters) >= min_clusters:
            return clusters

        
    return []

###
### Расстояние между двумя точками 
###
### point1 - Первая точка
### point2 - Вторая точка 
###
def distance(point1, point2):
    dimensions = len(point1)
    result = 0
    for i in range(0, len(point1)):
        result += (point1[i] - point2[i]) ** 2

    result = result ** 0.5
    return result

###
### Получение всех точек из массива, которые находятся в гиперсфере 
### указанного радиуса из указанным центром
###
### points - Массив точек
### center - Центр гиперсферы
### radius - Радиус гиперсферы
###
def get_hypersphere(points, center, radius):
    result = []
    indexes = []
    i = 0
    for point in points:
        if (distance(point, center) < radius):
            result.append(point)
            indexes.append(i)

        i = i + 1

    return result

###
### Нахождение центра масс набора точек 
### 
### points - Набор точек 
###
def center_list(points):
    dimensions = len(points[0])
    result = [0 for _ in range(dimensions)]
    for i in range(dimensions):
        for point in points:
            result[i] += point[i]

        result[i] /= len(points)

    return result

###
### Нахождение минимума из набора чисел
### 
### list - Набор чисел
###
def min_list(list):
    result = list[0]
    for elem in list:
        if elem < result:
            result = elem

    return result

###
### Нахождение максимума из набора чисел
### 
### list - Набор чисел
###
def max_list(list):
    result = list[0]
    for elem in list:
        if elem > result:
            result = elem

    return result

###
### Точка входа в программу 
###
if __name__ == "__main__":
    points = []
    min_clusters = 0

    with open("input.csv") as file:
        lines = file.readlines()
        
        line_count = 0
        for line in lines:
            arr = [int(elem) for elem in line.split(",")]
            if line_count == 0:
                min_clusters = arr[0]
            else:
                points.append(arr)

            line_count += 1

    with open("output.csv", "w") as file:
        cluster_index = 0
        clusters = forel(points, min_clusters)
        for cluster in clusters:
            for point in cluster:
                point.append(cluster_index)
                str_point = [str(elem) for elem in point]
                file.write("%s\n" % ",".join(str_point))

            cluster_index += 1

    
    
        
        
        
