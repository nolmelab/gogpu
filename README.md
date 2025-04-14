# WebGPU

Go bindings for WebGPU, a cross-platform, safe graphics API. It runs natively using [wgpu-native](https://github.com/gfx-rs/wgpu-native) on Vulkan, Metal, D3D12, and OpenGL ES based on https://github.com/rajveermalviya/go-webgpu. It also comes with web (JS) support based on https://github.com/mokiat/wasmgpu.

For more information, see:

- [WebGPU](https://gpuweb.github.io/gpuweb/)
- [WGSL](https://gpuweb.github.io/gpuweb/wgsl/)
- [webgpu-native](https://github.com/webgpu-native/webgpu-headers)

The included static libraries are built via [GitHub Actions](.github/workflows/build-wgpu.yml).

## Examples

|[boids][b]|[cube][c]|[triangle][t]|
:-:|:-:|:-:
| [![b-i]][b] | [![c-i]][c] | [![t-i]][t] |

[b-i]: https://raw.githubusercontent.com/rajveermalviya/go-webgpu/main/tests/boids/image-msaa.png
[b]: examples/boids
[c-i]: https://raw.githubusercontent.com/rajveermalviya/go-webgpu/main/tests/cube/image-msaa.png
[c]: examples/cube
[t-i]: https://raw.githubusercontent.com/rajveermalviya/go-webgpu/main/tests/triangle/image-msaa.png
[t]: examples/triangle


## 학습 

webgpu를 사용하여 go로 렌더러와 CG(계산 기하), 물리(충돌 중심) 구현을 연습하고자 합니다. 

CG와 물리는 게임 서버에서도 많이 사용되는 핵심 기능입니다. 짬짬이 문서를 작성하면서 상세하게 
이해할 수 있는 자료와 코드를 축적하고자 합니다. 

