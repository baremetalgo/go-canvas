package core

type Canvas struct {
	BaseLayer       *Layer
	Layers          []*Layer
	Brush           *Brush
	Active_Layer_Id uint32
}

func NewCanvas() *Canvas {
	base_layer := NewLayer()
	brush := NewBrush()

	canvas := Canvas{}
	canvas.BaseLayer = base_layer
	canvas.Layers = append(canvas.Layers, base_layer)
	canvas.Active_Layer_Id = base_layer.Id
	canvas.Brush = brush
	return &canvas
}

func (c *Canvas) GetActiveLayer() *Layer {
	active_layer := c.BaseLayer
	for _, layer := range c.Layers {
		if layer.Id == c.Active_Layer_Id {
			active_layer = layer
		}
	}
	return active_layer
}
