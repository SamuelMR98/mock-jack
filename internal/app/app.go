package app

const (
	screenWidth  = 960
	screenHeight = 540
)

type App struct {
}

func New() *App {
	a := &App{}
	return a
}

func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}