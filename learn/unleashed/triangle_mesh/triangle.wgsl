struct VertexOutput {
    @builtin(position) clip_position: vec4<f32>,
    @location(0) color: vec4<f32>
}

@vertex 
fn vs_main(@location(0) inPos: vec3<f32>) -> VertexOutput {
    var out: VertexOutput; 
    out.clip_position = vec4<f32>(inPos, 1.0);
    out.color = vec4<f32>(0, 0, 1, 1);
    return out;
}

@fragment 
fn fs_main(in : VertexOutput) -> @location(0) vec4<f32> {
    return in.color;
}

// 실수들: 
// - vec4를 vec으로 오타. 
// - VertexOutput 오타 
// 
// wgsl이 컴파일 할 때 오류 위치와 내용을 매우 자세히 잘 알려준다. 
// 