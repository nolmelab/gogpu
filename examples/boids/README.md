# boids 예제 

## 흐름 

- glfw로 창을 만든다. 
- wgpu의 Adapter 생성 
- wgpu의 Device 생성 
- Device의 큐 얻기 
- 렌더링 Surface 만들고 붙이기 
- 쉐이더 만들기 
  - Compute Shader 
  - Draw Shader 
- 렌더링 파이프라인 만들기 

렌더링 파이프라인 만들기가 핵심이다. 각 단계별로 여러 다양한 옵션이 있다. 
렌더링의 특징이 여러 옵션이 있다는 점이고 세부적인 제어가 필요하다는 것이다. 
이를 제어하는 것이 렌더링의 즐거움이다. 

## 렌더링 파이프라인 만들기 

s.device.CreateRenderPipeline() 

아래는 grok에 질문한 내용을 정리한 내용이다. 

WebGPU API에서 렌더링 파이프라인을 생성하는 메서드로, 그래픽스 렌더링(예: 버텍스 처리, 프래그먼트 셰이딩 등)을 
정의합니다. 이전 대화에서 다룬 Boids 예제와 cogentcore/webgpu를 기반으로, CreateRenderPipeline의 역할과 
Boids 예제에서의 사용을 자세히 설명하겠습니다. cogentcore/webgpu는 wgpu-native를 기반으로 하므로, 
해당 API를 기준으로 설명하며, s.device는 wgpu.Device 객체를 의미한다고 가정합니다.

1. CreateRenderPipeline란?
CreateRenderPipeline은 WebGPU에서 그래픽스 파이프라인을 생성하는 핵심 메서드입니다. 렌더링 파이프라인은 
버텍스 데이터(예: 위치, 색상)를 받아 화면에 픽셀을 그리는 과정을 정의합니다. 주요 구성 요소는 다음과 같습니다:

- 버텍스 셰이더: 3D 좌표를 화면 좌표로 변환.
- 프래그먼트 셰이더: 각 픽셀의 색상 결정.
- 프리미티브 설정: 점, 선, 삼각형 등의 렌더링 방식.
- 렌더 타겟 포맷: 출력 텍스처의 색상 포맷(예: BGRA8Unorm).
- 바인드 그룹 레이아웃: 셰이더에서 사용할 리소스(버퍼, 텍스처 등).

cogentcore/webgpu에서는 wgpu.Device.CreateRenderPipeline 메서드로 호출되며, 
wgpu.RenderPipelineDescriptor를 인자로 받아 파이프라인을 구성합니다.

2. 메서드 시그니처
cogentcore/webgpu에서 CreateRenderPipeline의 Go 코드 시그니처는 다음과 같습니다:

```go
func (d *Device) CreateRenderPipeline(descriptor *RenderPipelineDescriptor) (*RenderPipeline, error)
```
입력:
    descriptor: wgpu.RenderPipelineDescriptor 구조체로, 파이프라인의 모든 설정 포함.
출력: 
    *RenderPipeline: 생성된 렌더링 파이프라인 객체.
error: 생성 실패 시 오류(예: 셰이더 오류, 포맷 불일치).

RenderPipelineDescriptor의 주요 필드는:
- Layout: 바인드 그룹 레이아웃(리소스 바인딩 정의).
- Vertex: 버텍스 셰이더와 버퍼 레이아웃.
- Fragment: 프래그먼트 셰이더와 출력 포맷.
- Primitive: 프리미티브 타입(점, 선, 삼각형 등).
- DepthStencil: 깊이/스텐실 설정(선택적).
- Multisample: 멀티샘플링 설정.

3. Boids 예제에서의 CreateRenderPipeline

Boids 예제에서는 1000개의 Boid를 점(Point)으로 렌더링하기 위해 CreateRenderPipeline을 사용합니다. 
각 Boid의 위치는 GPU 버퍼에 저장되며, 버텍스 셰이더가 이를 화면 좌표로 변환하고, 
프래그먼트 셰이더가 주황색으로 색상을 지정합니다.

관련 코드 (Boids 예제)
```go
renderShader := `
    struct Boid {
        pos: vec2<f32>,
        vel: vec2<f32>,
    };
    struct VertexOutput {
        @builtin(position) pos: vec4<f32>,
    };
    @group(0) @binding(0) var<storage, read> boids: array<Boid>;
    @vertex
    fn vs_main(@builtin(vertex_index) idx: u32) -> VertexOutput {
        var boid = boids[idx];
        var pos = boid.pos;
        var output: VertexOutput;
        output.pos = vec4<f32>((pos.x / 400.0 - 1.0), -(pos.y / 300.0 - 1.0), 0.0, 1.0);
        return output;
    }
    @fragment
    fn fs_main() -> @location(0) vec4<f32> {
        return vec4<f32>(1.0, 0.5, 0.2, 1.0);
    }
`

renderShaderModule, err := device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
    WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{Code: renderShader},
})
if err != nil {
    panic(err)
}
defer renderShaderModule.Release()

renderBindGroupLayout, err := device.CreateBindGroupLayout(&wgpu.BindGroupLayoutDescriptor{
    Entries: []wgpu.BindGroupLayoutEntry{
        {
            Binding:    0,
            Visibility: wgpu.ShaderStageVertex,
            Buffer: wgpu.BufferBindingLayout{
                Type:           wgpu.BufferBindingTypeStorage,
                MinBindingSize: boidBufferSize,
            },
        },
    },
})
if err != nil {
    panic(err)
}
defer renderBindGroupLayout.Release()

renderPipelineLayout, err := device.CreatePipelineLayout(&wgpu.PipelineLayoutDescriptor{
    BindGroupLayouts: []wgpu.BindGroupLayout{renderBindGroupLayout},
})
if err != nil {
    panic(err)
}
defer renderPipelineLayout.Release()

renderPipeline, err := device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
    Layout: renderPipelineLayout,
    Vertex: wgpu.VertexState{
        Module:     renderShaderModule,
        EntryPoint: "vs_main",
    },
    Fragment: &wgpu.FragmentState{
        Module:     renderShaderModule,
        EntryPoint: "fs_main",
        Targets: []wgpu.ColorTargetState{
            {Format: wgpu.TextureFormatBGRA8Unorm},
        },
    },
    Primitive: wgpu.PrimitiveState{
        Topology: wgpu.PrimitiveTopologyPointList,
    },
})
if err != nil {
    panic(err)
}
defer renderPipeline.Release()
```

4. 코드 분석

CreateRenderPipeline 호출과 관련된 각 구성 요소를 단계별로 설명합니다.

(1) 셰이더 모듈 생성

```go

renderShaderModule, err := device.CreateShaderModule(&wgpu.ShaderModuleDescriptor{
    WGSLDescriptor: &wgpu.ShaderModuleWGSLDescriptor{Code: renderShader},
})
```

설명:
     WGSL(WebGPU Shading Language)로 작성된 셰이더 코드를 컴파일해 셰이더 모듈 생성.
    - vs_main: 버텍스 셰이더, Boid 위치를 화면 좌표로 변환.
    - fs_main: 프래그먼트 셰이더, 고정된 주황색(RGBA: 1.0, 0.5, 0.2, 1.0) 출력.
    - renderShaderModule은 버텍스와 프래그먼트 셰이더를 모두 포함.

(2) 바인드 그룹 레이아웃

```go
renderBindGroupLayout, err := device.CreateBindGroupLayout(&wgpu.BindGroupLayoutDescriptor{
    Entries: []wgpu.BindGroupLayoutEntry{
        {
            Binding:    0,
            Visibility: wgpu.ShaderStageVertex,
            Buffer: wgpu.BufferBindingLayout{
                Type:           wgpu.BufferBindingTypeStorage,
                MinBindingSize: boidBufferSize,
            },
        },
    },
})
```

설명:
셰이더에서 사용할 리소스(Boid 데이터 버퍼)의 레이아웃 정의.

- Binding: 0: 셰이더의 @group(0) @binding(0)에 매핑.
- Visibility: wgpu.ShaderStageVertex: 버텍스 셰이더에서만 접근.
- BufferBindingLayout.Type: wgpu.BufferBindingTypeStorage: Boid 데이터를 스토리지 버퍼로 처리.
- MinBindingSize: 버퍼 크기(1000 Boids × 16바이트) 지정.

(3) 파이프라인 레이아웃

```go
renderPipelineLayout, err := device.CreatePipelineLayout(&wgpu.PipelineLayoutDescriptor{
    BindGroupLayouts: []wgpu.BindGroupLayout{renderBindGroupLayout},
})
```
설명:
- 파이프라인에서 사용할 모든 바인드 그룹 레이아웃을 결합.
- Boids 예제는 단일 바인드 그룹(Boid 버퍼)만 사용하므로 하나만 포함.
- renderPipelineLayout은 CreateRenderPipeline의 Layout 필드에 전달됨.

(4) CreateRenderPipeline 호출
```go
renderPipeline, err := device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{
    Layout: renderPipelineLayout,
    Vertex: wgpu.VertexState{
        Module:     renderShaderModule,
        EntryPoint: "vs_main",
    },
    Fragment: &wgpu.FragmentState{
        Module:     renderShaderModule,
        EntryPoint: "fs_main",
        Targets: []wgpu.ColorTargetState{
            {Format: wgpu.TextureFormatBGRA8Unorm},
        },
    },
    Primitive: wgpu.PrimitiveState{
        Topology: wgpu.PrimitiveTopologyPointList,
    },
})
```

구성 요소:
- Layout: renderPipelineLayout:
- 셰이더 리소스 바인딩 정의.
- Vertex:
    - Module: renderShaderModule: WGSL 셰이더 모듈.
    - EntryPoint: "vs_main": 버텍스 셰이더 진입점.
    - 역할: 각 Boid의 위치를 읽어 화면 좌표(-1~1)로 변환.
    - 참고: VertexBufferLayout은 생략됨(직접 vertex_index로 Boid 배열 접근).
- Fragment:
    - Module: renderShaderModule, EntryPoint: "fs_main": 프래그먼트 셰이더.
    - Targets: 출력 텍스처 포맷(BGRA8Unorm) 지정. 스왑체인의 포맷과 일치해야 함.
    - 역할: 모든 Boid를 주황색으로 렌더링.
- Primitive:
    - Topology: wgpu.PrimitiveTopologyPointList: 각 Boid를 점으로 렌더링.
    - 대안: TriangleList로 삼각형 렌더링 가능(추가 버텍스 데이터 필요).
- 생략된 필드:
    - DepthStencil: 깊이 테스트 불필요(2D 렌더링).
    - Multisample: 기본값(1 샘플, 안티앨리어싱 없음).

(5) 렌더링 패스에서의 사용
생성된 renderPipeline은 렌더링 패스에서 설정됨:

```go
renderPass := commandEncoder.BeginRenderPass(&wgpu.RenderPassDescriptor{
    ColorAttachments: []wgpu.RenderPassColorAttachment{
        {
            View:       view,
            LoadOp:     wgpu.LoadOpClear,
            StoreOp:    wgpu.StoreOpStore,
            ClearValue: wgpu.Color{R: 0, G: 0, B: 0, A: 1},
        },
    },
})
renderPass.SetPipeline(renderPipeline)
renderPass.SetBindGroup(0, renderBindGroup, nil)
renderPass.Draw(numBoids, 1, 0, 0)
renderPass.End()
```

설명:
    - SetPipeline(renderPipeline): 이 파이프라인을 렌더링에 사용.
    - SetBindGroup: Boid 데이터 버퍼를 셰이더에 바인딩.
    - Draw(numBoids, 1, 0, 0): 1000개의 점을 그림(인스턴스 1개, 오프셋 0).

5. cogentcore/webgpu와의 호환성

API 일관성:
    - cogentcore/webgpu는 wgpu-native 기반이므로 CreateRenderPipeline의 동작은 표준 
      WebGPU와 동일.
    - wgpu.RenderPipelineDescriptor 구조체 필드와 플래그(ShaderStageVertex, 
      PrimitiveTopologyPointList 등)는 wgpu.h에서 파생됨.
주의점:
    - cogentcore/webgpu가 Cogent Core 프레임워크와 통합되어 있다면, 추가 유틸리티(예: 셰이더 로더)가 
      사용될 수 있음. 소스 확인 필요.
    - 메서드명이나 에러 처리 방식이 약간 다를 경우, github.com/cogentcore/webgpu의 API 문서 참고.

확인 방법:
    - device.CreateRenderPipeline 호출 전후에 디버그 로그 추가:

```go
fmt.Println("Creating render pipeline...")
renderPipeline, err := device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{...})
if err != nil {
    fmt.Printf("Error: %v\n", err)
}
```

6. 확장 및 개선

Boids 예제의 렌더링 파이프라인을 확장할 수 있는 방법:

삼각형 렌더링:
PrimitiveTopologyPointList → PrimitiveTopologyTriangleList.

각 Boid에 3개의 버텍스 정의:

```wgsl

struct Vertex {
    pos: vec3<f32>,
};
@group(0) @binding(0) var<storage, read> vertices: array<Vertex>;
```

VertexBufferLayout 추가:

```go
Vertex: wgpu.VertexState{
    Module:     renderShaderModule,
    EntryPoint: "vs_main",
    Buffers: []wg CarmelVertexBufferLayout{
        {
            ArrayStride: 12, // 3 floats
            Attributes: []wgpu.VertexAttribute{
                {Format: wgpu.VertexFormatFloat32x3, Offset: 0, ShaderLocation: 0},
            },
        },
    },
},
```
텍스처 매핑:
    - Boid를 점 대신 텍스처(예: 새 이미지)로 렌더링.
    - TextureUsageTextureBinding 텍스처와 샘플러 추가:

```go
texture, _ := device.CreateTexture(&wgpu.TextureDescriptor{...})
sampler, _ := device.CreateSampler(&wgpu.SamplerDescriptor{...})
```
셰이더 수정:
```wgsl
@group(0) @binding(1) var tex: texture_2d<f32>;
@group(0) @binding(2) var samp: sampler;
```
다중 파이프라인:
    - 그림자, 윤곽선 등 추가 효과를 위해 별도의 파이프라인 생성:

```go
outlinePipeline, _ := device.CreateRenderPipeline(&wgpu.RenderPipelineDescriptor{...})
```
깊이 테스트:
    - 3D Boids로 확장 시 깊이 버퍼 추가:

```go
DepthStencil: &wgpu.DepthStencilState{
    Format: wgpu.TextureFormatDepth24Plus,
    DepthWriteEnabled: true,
    DepthCompare: wgpu.CompareFunctionLess,
},
```

7. 문제 해결 팁

셰이더 오류:
    - CreateRenderPipeline 실패 시 WGSL 코드 문법 확인:

```text
[wgpu error]: Shader compilation failed: ...
```

device.CreateShaderModule 단계에서 오류 로그 활성화:
```go
shaderModule, err := device.CreateShaderModule(...)
if err != nil {
    fmt.Printf("Shader error: %v\n", err)
}
```

포맷 불일치:
    - Fragment.Targets의 Format이 스왑체인 포맷(BGRA8Unorm)과 일치해야 함.
    - 불일치 시:
        [wgpu error]: Render pipeline output format does not match attachment
리소스 바인딩 오류:
    - Layout과 BindGroupLayout이 셰이더의 @group, @binding과 매핑되지 않으면 런타임 오류.
    - 바인드 그룹 디버깅:
```go
fmt.Printf("BindGroupLayout entries: %+v\n", renderBindGroupLayout.Entries)
```
메모리 관리:
    renderPipeline, shaderModule, pipelineLayout 등은 Release()로 해제:
```go
defer renderPipeline.Release()
```

8. 결론

device.CreateRenderPipeline은 WebGPU에서 렌더링 파이프라인을 생성하는 핵심 메서드로, 
Boids 예제에서는 1000개의 Boid를 점으로 렌더링하기 위해 사용됩니다. 
RenderPipelineDescriptor를 통해 셰이더, 프리미티브 타입, 출력 포맷 등을 설정하며, 
cogentcore/webgpu에서도 wgpu-native 기반으로 동일하게 사용합니다. 