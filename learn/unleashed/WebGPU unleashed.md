# WebGPU Unleashed 

[webgpu unleashed](https://shi-yan.github.io/webgpuunleashed)


## GPU Driver 

- user mode driver / kernel mode driver 
  - WebGPU / (OpenGL, DirectX11/12, Vulkan, ...)

3 roles: 
- compiles api requests into machine code that GPU understands 
- resource manager 
- scheduler 
  - schedule and coordinate(synchronize) 

기능은 단순하다. 명령을 드라이버로 보내고 필요한 자원을 함께 제공하면 드라이버는 
쉐이더로 계산하여 그린다. 그것이 근간이고 이 위에 다양한 알고리즘을 만든다. 
그래서, 예쁜 세상을 표한한다. 

## The GPU Pipeline 

픽셀 공장. 

- pipeline configruations 
  - shader programs 
  - and lots other 
- feed the pipeline with inputs 
  - millions of triangles and (not too many) texures
- other data 
  - uniform buffers 
  - they are arguments to the shader

- triangles to pixels: 
    - vertex 
    - rasterization 
    - fragment 
    - blending 

## Creating an Empty Canvas 

javascript로 web에서 튜토리얼을 진행한다. 이를 glfw와 go로 변경하여 따라간다. 

그릴 Surface, Adapter를 얻는다. 
canvas config는 surface config와 같다. 

- instance 
- adapter 
- surface
- device 
  - queue 
- render pass 

colorAttachment는 컬러 버퍼에 해당한다.  surface에서 텍스처를 얻고 다시 view를 얻어서 연결한다. 

## Drawing a triangle 

unleashed 폴더에 각 챕터에 해당하는 폴더를 추가하고 코드에 직접 설명을 적는다. 


## 자료 

- [WebGPU Rocks](https://webgpu.rocks/) 
  - API와 WGSL에 대한 간략한 설명을 포함한다
  - js API이다 
  - 기억할 필요 없이 사이트만 기억하면 된다
  - 문서가 깔끔하다







