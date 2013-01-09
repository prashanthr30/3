package ptx

//This file is auto-generated. Editing is futile.

func init() { Code["reducemaxdiff"] = REDUCEMAXDIFF }

const REDUCEMAXDIFF = `
//
// Generated by NVIDIA NVVM Compiler
// Compiler built on Sat Sep 22 02:35:14 2012 (1348274114)
// Cuda compilation tools, release 5.0, V0.2.1221
//

.version 3.1
.target sm_30
.address_size 64

	.file	1 "/tmp/tmpxft_00006be9_00000000-9_reducemaxdiff.cpp3.i"
	.file	2 "/home/arne/src/code.google.com/p/mx3/gpu/ptx/reducemaxdiff.cu"
	.file	3 "/usr/local/cuda-5.0/nvvm/ci_include.h"
	.file	4 "/usr/local/cuda/bin/../include/sm_11_atomic_functions.h"
// __cuda_local_var_33851_32_non_const_sdata has been demoted

.visible .entry reducemaxdiff(
	.param .u64 reducemaxdiff_param_0,
	.param .u64 reducemaxdiff_param_1,
	.param .u64 reducemaxdiff_param_2,
	.param .f32 reducemaxdiff_param_3,
	.param .u32 reducemaxdiff_param_4
)
{
	.reg .pred 	%p<8>;
	.reg .s32 	%r<41>;
	.reg .f32 	%f<31>;
	.reg .s64 	%rd<16>;
	// demoted variable
	.shared .align 4 .b8 __cuda_local_var_33851_32_non_const_sdata[2048];

	ld.param.u64 	%rd5, [reducemaxdiff_param_0];
	ld.param.u64 	%rd6, [reducemaxdiff_param_1];
	ld.param.u64 	%rd7, [reducemaxdiff_param_2];
	ld.param.f32 	%f30, [reducemaxdiff_param_3];
	ld.param.u32 	%r9, [reducemaxdiff_param_4];
	cvta.to.global.u64 	%rd1, %rd7;
	cvta.to.global.u64 	%rd2, %rd6;
	cvta.to.global.u64 	%rd3, %rd5;
	.loc 2 8 1
	mov.u32 	%r40, %ntid.x;
	mov.u32 	%r10, %ctaid.x;
	mov.u32 	%r2, %tid.x;
	mad.lo.s32 	%r39, %r40, %r10, %r2;
	mov.u32 	%r11, %nctaid.x;
	mul.lo.s32 	%r4, %r40, %r11;
	.loc 2 8 1
	setp.ge.s32 	%p1, %r39, %r9;
	@%p1 bra 	BB0_2;

BB0_1:
	.loc 2 8 1
	mul.wide.s32 	%rd8, %r39, 4;
	add.s64 	%rd9, %rd3, %rd8;
	add.s64 	%rd10, %rd2, %rd8;
	ld.global.f32 	%f5, [%rd10];
	ld.global.f32 	%f6, [%rd9];
	sub.f32 	%f7, %f6, %f5;
	.loc 3 395 5
	abs.f32 	%f8, %f7;
	.loc 3 435 5
	max.f32 	%f30, %f30, %f8;
	.loc 2 8 1
	add.s32 	%r39, %r39, %r4;
	.loc 2 8 1
	setp.lt.s32 	%p2, %r39, %r9;
	@%p2 bra 	BB0_1;

BB0_2:
	.loc 2 8 1
	mul.wide.s32 	%rd11, %r2, 4;
	mov.u64 	%rd12, __cuda_local_var_33851_32_non_const_sdata;
	add.s64 	%rd4, %rd12, %rd11;
	st.shared.f32 	[%rd4], %f30;
	bar.sync 	0;
	.loc 2 8 1
	setp.lt.u32 	%p3, %r40, 66;
	@%p3 bra 	BB0_6;

BB0_3:
	.loc 2 8 1
	mov.u32 	%r7, %r40;
	shr.u32 	%r40, %r7, 1;
	.loc 2 8 1
	setp.ge.u32 	%p4, %r2, %r40;
	@%p4 bra 	BB0_5;

	.loc 2 8 1
	ld.shared.f32 	%f9, [%rd4];
	add.s32 	%r16, %r40, %r2;
	mul.wide.u32 	%rd13, %r16, 4;
	add.s64 	%rd15, %rd12, %rd13;
	ld.shared.f32 	%f10, [%rd15];
	.loc 3 435 5
	max.f32 	%f11, %f9, %f10;
	.loc 2 8 1
	st.shared.f32 	[%rd4], %f11;

BB0_5:
	.loc 2 8 1
	bar.sync 	0;
	.loc 2 8 1
	setp.gt.u32 	%p5, %r7, 131;
	@%p5 bra 	BB0_3;

BB0_6:
	.loc 2 8 1
	setp.gt.s32 	%p6, %r2, 31;
	@%p6 bra 	BB0_8;

	.loc 2 8 1
	ld.volatile.shared.f32 	%f12, [%rd4];
	ld.volatile.shared.f32 	%f13, [%rd4+128];
	.loc 3 435 5
	max.f32 	%f14, %f12, %f13;
	.loc 2 8 1
	st.volatile.shared.f32 	[%rd4], %f14;
	ld.volatile.shared.f32 	%f15, [%rd4+64];
	ld.volatile.shared.f32 	%f16, [%rd4];
	.loc 3 435 5
	max.f32 	%f17, %f16, %f15;
	.loc 2 8 1
	st.volatile.shared.f32 	[%rd4], %f17;
	ld.volatile.shared.f32 	%f18, [%rd4+32];
	ld.volatile.shared.f32 	%f19, [%rd4];
	.loc 3 435 5
	max.f32 	%f20, %f19, %f18;
	.loc 2 8 1
	st.volatile.shared.f32 	[%rd4], %f20;
	ld.volatile.shared.f32 	%f21, [%rd4+16];
	ld.volatile.shared.f32 	%f22, [%rd4];
	.loc 3 435 5
	max.f32 	%f23, %f22, %f21;
	.loc 2 8 1
	st.volatile.shared.f32 	[%rd4], %f23;
	ld.volatile.shared.f32 	%f24, [%rd4+8];
	ld.volatile.shared.f32 	%f25, [%rd4];
	.loc 3 435 5
	max.f32 	%f26, %f25, %f24;
	.loc 2 8 1
	st.volatile.shared.f32 	[%rd4], %f26;
	ld.volatile.shared.f32 	%f27, [%rd4+4];
	ld.volatile.shared.f32 	%f28, [%rd4];
	.loc 3 435 5
	max.f32 	%f29, %f28, %f27;
	.loc 2 8 1
	st.volatile.shared.f32 	[%rd4], %f29;

BB0_8:
	.loc 2 8 1
	setp.ne.s32 	%p7, %r2, 0;
	@%p7 bra 	BB0_10;

	.loc 2 8 1
	ld.shared.u32 	%r37, [__cuda_local_var_33851_32_non_const_sdata];
	.loc 3 1881 5
	atom.global.max.s32 	%r38, [%rd1], %r37;

BB0_10:
	.loc 2 9 2
	ret;
}


`
