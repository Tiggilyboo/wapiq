package main

type Query struct {
	Request Request
	Map     Map
	Args    MapResult
}

func (q *Query) Invoke() ([]MapResult, error) {
	var vs string
	var e bool
	for k, v := range q.Args {
		vs = v.(string)
		_, e = q.Request.Query[k]
		if e {
			q.Request.Query.Set(k, vs)
		}
		_, e = q.Request.Head[k]
		if e {
			q.Request.Head.Set(k, vs)
		}
		_, e = q.Request.Body[k]
		if e {
			q.Request.Body.Set(k, vs)
		}
	}
	j, err := q.Map.Invoke(&q.Request)
	if err != nil {
		return nil, err
	}
	return j, nil
}
