//go:build !js

package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"os"
	"runtime"
	"unsafe"

	"github.com/cogentcore/webgpu/wgpu"
	"github.com/cogentcore/webgpu/wgpuglfw"
	"github.com/go-gl/glfw/v3.3/glfw"

	_ "embed"
)

func main() {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ClientAPI, glfw.NoAPI)
	window, err := glfw.CreateWindow(640, 480, "go-webgpu with glfw", nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	s, err := InitState(window, wgpuglfw.GetSurfaceDescriptor(window))
	if err != nil {
		panic(err)
	}
	defer s.Destroy()

	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		// Print resource usage on pressing 'R'
		if key == glfw.KeyR && (action == glfw.Press || action == glfw.Repeat) {
			report := s.instance.GenerateReport()
			buf, _ := json.MarshalIndent(report, "", "  ")
			fmt.Print(string(buf))
		}
	})

	window.SetSizeCallback(func(_ *glfw.Window, width, height int) {
		s.Resize(width, height)
	})

	for !window.ShouldClose() {
		glfw.PollEvents()

		err := s.Render()
		if err != nil {
			fmt.Println("error occurred while rendering:", err)

			errstr := err.Error()
			switch {
			case strings.Contains(errstr, "Surface timed out"): // do nothing
			case strings.Contains(errstr, "Surface is outdated"): // do nothing
			case strings.Contains(errstr, "Surface was lost"): // do nothing
			default:
				panic(err)
			}
		}
	}
}

var forceFallbackAdapter = os.Getenv("WGPU_FORCE_FALLBACK_ADAPTER") == "1"

func init() {
	runtime.LockOSThread()

	switch os.Getenv("WGPU_LOG_LEVEL") {
	case "OFF":
		wgpu.SetLogLevel(wgpu.LogLevelOff)
	case "ERROR":
		wgpu.SetLogLevel(wgpu.LogLevelError)
	case "WARN":
		wgpu.SetLogLevel(wgpu.LogLevelWarn)
	case "INFO":
		wgpu.SetLogLevel(wgpu.LogLevelInfo)
	case "DEBUG":
		wgpu.SetLogLevel(wgpu.LogLevelDebug)
	case "TRACE":
		wgpu.SetLogLevel(wgpu.LogLevelTrace)
	}
}

type State struct {
	instance  *wgpu.Instance
	adapter   *wgpu.Adapter
	surface   *wgpu.Surface
	device    *wgpu.Device
	queue     *wgpu.Queue
	config    *wgpu.SurfaceConfiguration
	pipeline  *wgpu.RenderPipeline
	vertexBuf *wgpu.Buffer
}

// embed annotation tells go to embed the file into shader string

//go:embed triangle.wgsl
var shader string

type Vertex struct {
	pos   [3]float32
	color [4]float32
}

var VertexBufferLayout = wgpu.VertexBufferLayout{
	ArrayStride: uint64(unsafe.Sizeof(Vertex{})),
	StepMode:    wgpu.VertexStepModeVertex,
	Attributes: []wgpu.VertexAttribute{
		{
			Format:         wgpu.VertexFormatFloat32x3,
			Offset:         0,
			ShaderLocation: 0,
		},
		{
			Format:         wgpu.VertexFormatFloat32x4,
			Offset:         4 * 3,
			ShaderLocation: 1,
		},
	},
}

func createVertex(p1, p2, p3 float32, c1, c2, c3, alpha float32) Vertex {
	return Vertex{
		pos:   [3]float32{p1, p2, p3},
		color: [4]float32{c1, c2, c3, alpha},
	}
}

var vertexData = [...]Vertex{
	createVertex(0.0, 0.5, 0.0, 1.0, 0.0, 1.0, 1.0),
	createVertex(0.5, 0.0, 0.0, 1.0, 1.0, 0.0, 1.0),
	createVertex(1, 0.5, 0.0, 1.0, 1.0, 1.0, 1.0),

	createVertex(0.0, 0.5, 0.0, 1.0, 0.0, 1.0, 1.0),
	createVertex(-0.5, 0.0, 0.0, 1.0, 1.0, 1.0, 1.0),
	createVertex(0.5, 0.0, 0.0, 1.0, 1.0, 1.0, 1.0),

	createVertex(-0.5, 0.0, 0.0, 1.0, 0.0, 1.0, 1.0),
	createVertex(0.0, -0.5, 0.0, 1.0, 1.0, 1.0, 1.0),
	createVertex(0.5, 0.0, 0.0, 1.0, 1.0, 1.0, 1.0),

	createVertex(-0.5, 0.0, 0.0, 1.0, 0.0, 1.0, 1.0),
	createVertex(-1.0, -0.5, 0.0, 1.0, 0.0, 1.0, 1.0),
	createVertex(0.0, -0.5, 0.0, 1.0, 1.0, 0.0, 1.0),
}

func InitState[T interface{ GetSize() (int, int) }](window T, sd *wgpu.SurfaceDescriptor) (s *State, err error) {
	defer func() {
		if err != nil {
			s.Destroy()
			s = nil
		}
	}()
	s = &State{}

	s.instance = wgpu.CreateInstance(nil)

	s.surface = s.instance.CreateSurface(sd)

	s.adapter, err = s.instance.RequestAdapter(&wgpu.RequestAdapterOptions{
		ForceFallbackAdapter: forceFallbackAdapter,
		CompatibleSurface:    s.surface,
	})
	if err != nil {
		return s, err
	}
	defer s.adapter.Release()

	s.device, err = s.adapter.RequestDevice(nil)
	if err != nil {
		return s, err
	}
	s.queue = s.device.GetQueue()

	caps := s.surface.GetCapabilities(s.adapter)

	width, height := window.GetSize()
	s.config = &wgpu.SurfaceConfiguration{
		Usage:       wgpu.TextureUsageRenderAttachment,
		Format:      caps.Formats[0],
		Width:       uint32(width),
		Height:      uint32(height),
		PresentMode: wgpu.PresentModeFifo,
		AlphaMode:   caps.AlphaModes[0],
	}

	s.surface.Configure(s.adapter, s.device, s.config)

	shader, err := s.device.CreateShaderModule(
		&wgpu.ShaderModuleDescriptor{
			Label:          "traingle.wgsl",
			WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{Code: shader},
		},
	)

	if err != nil {
		return s, err
	}

	// create buffer
	s.vertexBuf, err = s.device.CreateBufferInit(&wgpu.BufferInitDescriptor{
		Label:    "Vertex Buffer",
		Contents: wgpu.ToBytes(vertexData[:]),
		Usage:    wgpu.BufferUsageVertex,
	})

	if err != nil {
		return s, err
	}

	// PipelineLayout, PipelineLayoutDesc

	s.pipeline, err = s.device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
		Label: "Render Pipeline",
		Vertex: wgpu.VertexState{
			Module:     shader,
			EntryPoint: "vs_main",
			Buffers:    []wgpu.VertexBufferLayout{VertexBufferLayout}, // 버텍스 버퍼 레이아웃 지정이 필수
		},
		Primitive: wgpu.PrimitiveState{
			Topology:         wgpu.PrimitiveTopologyTriangleList,
			StripIndexFormat: wgpu.IndexFormatUndefined,
			FrontFace:        wgpu.FrontFaceCCW,
			CullMode:         wgpu.CullModeNone,
		},
		Multisample: wgpu.MultisampleState{
			Count:                  1,
			Mask:                   0xFFFFFFFF,
			AlphaToCoverageEnabled: false,
		},
		Fragment: &wgpu.FragmentState{
			Module:     shader,
			EntryPoint: "fs_main",
			Targets: []wgpu.ColorTargetState{
				{
					Format:    s.config.Format,
					Blend:     &wgpu.BlendStateReplace,
					WriteMask: wgpu.ColorWriteMaskAll,
				},
			},
		},
	})

	return s, nil
}

func (s *State) Resize(width, height int) {
	if width > 0 && height > 0 {
		s.config.Width = uint32(width)
		s.config.Height = uint32(height)

		s.surface.Configure(s.adapter, s.device, s.config)
	}
}

func (s *State) Render() error {
	nextTexture, err := s.surface.GetCurrentTexture()
	if err != nil {
		return err
	}
	view, err := nextTexture.CreateView(nil)
	if err != nil {
		return err
	}
	defer view.Release()

	encoder, err := s.device.CreateCommandEncoder(&wgpu.CommandEncoderDescriptor{
		Label: "Command Encoder",
	})
	if err != nil {
		return err
	}
	defer encoder.Release()

	pass := encoder.BeginRenderPass(&wgpu.RenderPassDescriptor{
		ColorAttachments: []wgpu.RenderPassColorAttachment{
			{
				View:       view,
				LoadOp:     wgpu.LoadOpClear,
				StoreOp:    wgpu.StoreOpStore,
				ClearValue: wgpu.ColorGreen,
			},
		},
	})

	pass.SetPipeline(s.pipeline)
	pass.SetVertexBuffer(0, s.vertexBuf, 0, wgpu.WholeSize)
	pass.Draw(12, 4, 0, 0)
	pass.End()
	pass.Release() // must release

	cmdBuffer, err := encoder.Finish(nil)
	if err != nil {
		return err
	}
	defer cmdBuffer.Release()

	s.queue.Submit(cmdBuffer)
	s.surface.Present()

	return nil
}

func (s *State) Destroy() {
	if s.pipeline != nil {
		s.pipeline.Release()
		s.pipeline = nil
	}

	if s.config != nil {
		s.config = nil
	}

	if s.queue != nil {
		s.queue.Release()
		s.queue = nil
	}

	if s.device != nil {
		s.device.Release()
		s.device = nil
	}

	if s.surface != nil {
		s.surface.Release()
		s.surface = nil
	}

	if s.instance != nil {
		s.instance.Release()
		s.instance = nil
	}
}
