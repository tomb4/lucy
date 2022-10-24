package navmesh

import (
	"errors"
	detour "lucy/pkg/nav/Detour"
	dtcache "lucy/pkg/nav/DetourTileCache"
	"unsafe"
)

const (
	meshMaxNode = 2048
	pathMaxNode = 256

	NilPolyRef detour.DtPolyRef = 0
)

var (
	initNavMeshErr = errors.New("初始化寻路地图失败")
)

type NavMesh struct {
	dnm                 *detour.DtNavMesh
	dtc                 *dtcache.DtTileCache
	query               *detour.DtNavMeshQuery
	tileCacheNeedUpdate bool
}

type Point struct {
	X, Y, Z float32
}

func (p Point) ToF() []float32 {
	return (*[3]float32)(unsafe.Pointer(&p))[:]
}

func CoordinateTrans(p []float32) Point {
	return Point{X: p[0], Y: p[1], Z: p[2]}
}

func InitNavMesh(path string) (*NavMesh, error) {
	dnm, dtc, err := LoadDynamicMesh(path)
	if err != nil {
		return nil, err
	}
	m := &NavMesh{
		dnm: dnm,
		dtc: dtc,
	}
	query := detour.DtAllocNavMeshQuery()
	if query == nil {
		return nil, initNavMeshErr
	}
	stat := query.Init(dnm, meshMaxNode)
	if !detour.DtStatusSucceed(stat) {
		return nil, initNavMeshErr
	}
	m.query = query
	return m, nil
}

/*
*
根据起止点寻路
input: start,end 起止点的坐标
out: bool-寻路是否成功  []Point 从起点到终点途径的所有点位
*/
func (m *NavMesh) FindStraightPath(start, end Point) (bool, []Point) {
	startRef := m.findNearestPoly(start)
	if startRef == 0 {
		return false, nil
	}
	endRef := m.findNearestPoly(end)
	if endRef == 0 {
		return false, nil
	}

	path := make([]detour.DtPolyRef, pathMaxNode)
	pathCount := 0
	st := m.query.FindPath(startRef, endRef, start.ToF(), end.ToF(), detour.DtAllocDtQueryFilter(), path, &pathCount, pathMaxNode)
	if !detour.DtStatusSucceed(st) {
		return false, nil
	}

	if pathCount == 0 {
		return false, nil
	}

	sPath := [pathMaxNode * 3]float32{}
	sPathFlags := [pathMaxNode]detour.DtStraightPathFlags{}
	sPathPolys := [pathMaxNode]detour.DtPolyRef{}
	sPathCount := 0
	st = m.query.FindStraightPath(start.ToF(), end.ToF(), path, pathCount, sPath[:], sPathFlags[:], sPathPolys[:],
		&sPathCount, pathMaxNode, detour.DT_STRAIGHTPATH_AREA_CROSSINGS)
	if !detour.DtStatusSucceed(st) {
		return false, nil
	}

	ret := make([]Point, 0)
	for i := 0; i < sPathCount*3; {
		p := Point{}
		p.X = sPath[i]
		i++
		p.Y = sPath[i]
		i++
		p.Z = sPath[i]
		i++
		ret = append(ret, p)
	}
	return true, ret
}

/*
*
进行射线寻路
input: start,end 起止点的坐标
out: bool-寻路是否成功  Point 终点点位
*/
func (m *NavMesh) RayCast(start, end Point) (bool, Point) {
	startRef := m.findNearestPoly(start)
	hitPos := Point{}
	if startRef == 0 {
		return false, hitPos
	}
	endRef := m.findNearestPoly(end)
	if endRef == 0 {
		return false, hitPos
	}
	var t float32
	var hitNormal [3]float32
	var path [pathMaxNode]detour.DtPolyRef
	var pathCount int
	status := m.query.Raycast(startRef, start.ToF(), end.ToF(), detour.DtAllocDtQueryFilter(), &t, hitNormal[:], path[:], &pathCount, pathMaxNode)
	if !detour.DtStatusSucceed(status) {
		return false, hitPos
	}
	hit := (t <= 1)
	if hit {
		hitP := [3]float32{}
		detour.DtVlerp(hitP[:], start.ToF(), end.ToF(), t)
		if pathCount > 0 {
			success, height := m.GetPolyHeight(CoordinateTrans(hitP[:]))
			if !success {
				return false, hitPos
			}
			hitP[1] = height
		}
		hitPos = CoordinateTrans(hitP[:])
	} else {
		hitPos = end
		success, height := m.GetPolyHeight(end)
		if !success {
			return false, hitPos
		}
		hitPos.Y = height
	}
	return true, hitPos
}

/*
*
判断某个点是否可达
*/
func (m *NavMesh) CanReach(pos Point) bool {
	if m.findNearestPoly(pos) == NilPolyRef {
		return false
	}
	return true
}

/*
*
pos: 根据坐标点获取高度  Y值将被忽略
*/
func (m *NavMesh) GetPolyHeight(pos Point) (bool, float32) {
	ref := m.findNearestPoly(pos)
	if ref == NilPolyRef {
		return false, 0
	}

	var h float32
	st := m.query.GetPolyHeight(ref, pos.ToF(), &h)
	if !detour.DtStatusSucceed(st) {
		return false, 0
	}
	return true, h
}

/*
*
增加障碍物
*/
func (m *NavMesh) AddObstacle(pos Point, radius, height float32) (bool, dtcache.DtObstacleRef) {
	var ref dtcache.DtObstacleRef
	status := m.dtc.AddObstacle(pos.ToF(), radius, height, &ref)
	if !detour.DtStatusSucceed(status) {
		return false, 0
	}
	m.tileCacheNeedUpdate = true
	return true, ref
}

/*
*
删除障碍物
*/
func (m *NavMesh) RemoveObstacle(ref dtcache.DtObstacleRef) bool {
	status := m.dtc.RemoveObstacle(ref)
	if !detour.DtStatusSucceed(status) {
		return false
	}
	m.tileCacheNeedUpdate = true
	return true
}

/*
*
更新
*/
func (m *NavMesh) TileCacheUpdate() bool {
	if !m.tileCacheNeedUpdate {
		return true
	}
	status := m.dtc.Update(1, m.dnm, nil)
	if !detour.DtStatusSucceed(status) {
		return false
	}
	return true
}

// 根据点查询面的id
func (m *NavMesh) findNearestPoly(p Point) detour.DtPolyRef {
	pos := []float32{p.X, p.Y, p.Z}
	var ref detour.DtPolyRef
	extents := []float32{1, 100, 1}
	status := m.query.FindNearestPoly(pos, extents, detour.DtAllocDtQueryFilter(), &ref, pos)
	if !detour.DtStatusSucceed(status) {
		return NilPolyRef
	}
	return ref
}
