struct VertexOutput {
    @builtin(position) clip_position: vec4<f32>,
}

@vertex 
fn vs_main(@builtin(vertex_index) in_vertex_index: u32,) -> VertexOutput {
    var out: VertexOutput; 
    let x = f32(1 - i32(in_vertex_index)) * 0.5; 
    let y = f32(i32(in_vertex_index & 1u) * 2 - 1) * 0.5; 
    out.clip_position = vec4<f32>(x, y, 0.0, 1.0);
    return out;
}

@fragment 
fn fs_main(in : VertexOutput) -> @location(0) vec4<f32> {
    return vec4<f32>(0.3, 0.2, 0.1, 1.0);
}

// 실수들: 
// - vec4를 vec으로 오타. 
// - VertexOutput 오타 
// 
// wgsl이 컴파일 할 때 오류 위치와 내용을 매우 자세히 잘 알려준다. 
// 