package handler

import (
	"encoding/json"
	"strings"

	"github.com/hitpaw/mcp-server-hitpaw/internal/protocol"
)

// HandleListPhotoModels 列出所有图片增强模型
func (h *Handlers) HandleListPhotoModels(_ json.RawMessage) *protocol.CallToolResult {
	var sb strings.Builder
	sb.WriteString("📸 图片增强可用模型列表（共20个）\n")
	sb.WriteString("═══════════════════════════════════════\n\n")

	sb.WriteString("👤 人像清晰模型（柔和美颜效果）\n")
	sb.WriteString("───────────────────────────────────────\n")
	sb.WriteString("  face_2x\n")
	sb.WriteString("    放大: 2x | 人脸+背景双模型，柔和美颜\n")
	sb.WriteString("    触发: 含人脸 + \"清晰/修复/美颜/自拍/老照片\" + 2倍或未指定\n")
	sb.WriteString("    不适合: 纯风景 / 4倍放大 / 需保留皱纹胡须\n\n")
	sb.WriteString("  face_4x\n")
	sb.WriteString("    放大: 4x | 同上但4倍放大\n")
	sb.WriteString("    触发: 含人脸 + 明确要求\"4倍放大/大幅放大\"\n\n")

	sb.WriteString("👤 人像自然模型（保留真实纹理）\n")
	sb.WriteString("───────────────────────────────────────\n")
	sb.WriteString("  face_v2_2x\n")
	sb.WriteString("    放大: 2x | 保留皱纹/胡须/装饰等真实纹理\n")
	sb.WriteString("    触发: 含人脸 + \"保留纹理/自然/不磨皮\" + 2倍或未指定\n")
	sb.WriteString("    不适合: 柔和美颜 / 4倍放大\n\n")
	sb.WriteString("  face_v2_4x\n")
	sb.WriteString("    放大: 4x | 同上但4倍放大\n")
	sb.WriteString("    触发: 含人脸 + \"保留纹理\" + 明确4倍放大\n\n")

	sb.WriteString("🏞️ 通用高清模型（常规清晰化）\n")
	sb.WriteString("───────────────────────────────────────\n")
	sb.WriteString("  general_2x\n")
	sb.WriteString("    放大: 2x | 通用画面清晰化+细节锐化+画质优化\n")
	sb.WriteString("    触发: 无人脸 + \"清晰/画质提升\" + 2倍或未指定\n")
	sb.WriteString("    不适合: 含人脸 / 高保真 / 4倍放大\n\n")
	sb.WriteString("  general_4x\n")
	sb.WriteString("    放大: 4x | 同上但4倍放大\n\n")

	sb.WriteString("🎨 高保真模型（保留原图风格）\n")
	sb.WriteString("───────────────────────────────────────\n")
	sb.WriteString("  high_fidelity_2x\n")
	sb.WriteString("    放大: 2x | 消除伪影/噪声，保留原图纹理风格\n")
	sb.WriteString("    触发: \"高保真/保留风格/摄影作品/艺术照\" + 2倍\n")
	sb.WriteString("    适合: 摄影作品/艺术照/商品图/4K海报\n\n")
	sb.WriteString("  high_fidelity_4x\n")
	sb.WriteString("    放大: 4x | 同上但4倍放大\n\n")

	sb.WriteString("🔇 降噪模型（不放大，仅降噪优化）\n")
	sb.WriteString("───────────────────────────────────────\n")
	sb.WriteString("  sharpen_denoise_1x\n")
	sb.WriteString("    放大: 1x | 锐化降噪，去除噪点+锐化+生成新纹理\n")
	sb.WriteString("    触发: \"降噪/锐化/手机抓拍\" + 无放大需求\n")
	sb.WriteString("    不适合: 需放大 / 需保留原始纹理\n\n")
	sb.WriteString("  detail_denoise_1x\n")
	sb.WriteString("    放大: 1x | 极致保真降噪，去除复杂噪声+保留原纹理\n")
	sb.WriteString("    触发: \"降噪/保真/保留细节/夜景/复杂噪声\" + 无放大\n")
	sb.WriteString("    不适合: 需放大 / 需强锐化\n\n")

	sb.WriteString("🤖 生成式人像模型（SD，极低质量人像）\n")
	sb.WriteString("───────────────────────────────────────\n")
	sb.WriteString("  generative_portrait_1x\n")
	sb.WriteString("    放大: 1x | Diffusion人像增强，补全人脸五官\n")
	sb.WriteString("    触发: 极低质量人像(大头照/侧脸/老照片/表情包) + 无放大\n\n")
	sb.WriteString("  generative_portrait_2x\n")
	sb.WriteString("    放大: 2x | 同上+2倍放大\n\n")
	sb.WriteString("  generative_portrait_4x\n")
	sb.WriteString("    放大: 4x | 同上+4倍放大\n\n")

	sb.WriteString("🤖 生成式通用模型（SD，极低质量通用图）\n")
	sb.WriteString("───────────────────────────────────────\n")
	sb.WriteString("  generative_1x\n")
	sb.WriteString("    放大: 1x | Diffusion通用增强，生成自然纹理\n")
	sb.WriteString("    触发: 极低质量通用图 + \"生成补全/高清\" + 无放大\n\n")
	sb.WriteString("  generative_2x\n")
	sb.WriteString("    放大: 2x | 同上+2倍放大\n\n")
	sb.WriteString("  generative_4x\n")
	sb.WriteString("    放大: 4x | 同上+4倍放大\n\n")
	sb.WriteString("  generative_6x  ⭐ 新增\n")
	sb.WriteString("    放大: 6x | SD生成式6倍放大\n")
	sb.WriteString("    触发: 用户明确要求\"6倍放大\"（仅此模型支持6x）\n")
	sb.WriteString("    注意: SD输出上限34MP，请控制输入分辨率\n\n")
	sb.WriteString("  generative_8x  ⭐ 新增\n")
	sb.WriteString("    放大: 8x | SD生成式8倍放大\n")
	sb.WriteString("    触发: 用户明确要求\"8倍放大/极致放大\"（仅此模型支持8x）\n")
	sb.WriteString("    注意: SD输出上限34MP，请控制输入分辨率\n\n")

	sb.WriteString("⚡ 生成式快速模型（SD，效率优先）\n")
	sb.WriteString("───────────────────────────────────────\n")
	sb.WriteString("  generative_enhance_fast_2x\n")
	sb.WriteString("    放大: 2x | SD快速生成，效率优先，质量略低于generative_2x\n")
	sb.WriteString("    触发: 低质通用图 + \"快速/效率优先\" + 2倍放大\n\n")
	sb.WriteString("  generative_enhance_fast_4x\n")
	sb.WriteString("    放大: 4x | SD快速生成，效率优先，质量略低于generative_4x\n")
	sb.WriteString("    触发: 低质通用图 + \"快速/效率优先\" + 4倍放大\n\n")

	sb.WriteString("═══════════════════════════════════════\n")
	sb.WriteString("📌 限制信息：\n")
	sb.WriteString("  最大输入: 67MP | 最大输出: 432MP(非SD) / 34MP(SD)\n")
	sb.WriteString("  输入格式: .jpg .png .webp .bmp .tif .heic 等\n")
	sb.WriteString("  输出格式: .jpg .png .webp .bmp .tif 等\n\n")
	sb.WriteString("⚠️ 6x/8x 高倍放大说明：\n")
	sb.WriteString("  - 仅 generative_6x、generative_8x 两个模型支持\n")
	sb.WriteString("  - 其他模型系列（face/face_v2/general/high_fidelity）最高支持 4x\n")
	sb.WriteString("  - 6x/8x 为 SD 模型，输出上限 34MP，请注意输入图片尺寸\n")

	return protocol.SuccessResult(sb.String())
}

// HandleListVideoModels 列出所有视频增强模型
func (h *Handlers) HandleListVideoModels(_ json.RawMessage) *protocol.CallToolResult {
	var sb strings.Builder
	sb.WriteString("🎬 视频增强可用模型列表（共8个）\n")
	sb.WriteString("═══════════════════════════════════════\n\n")

	sb.WriteString("👤 人像/人脸模型\n")
	sb.WriteString("───────────────────────────────────────\n")
	sb.WriteString("  face_soft_2x         | 人脸柔和2x  | model_id=56\n")
	sb.WriteString("  portrait_restore_1x  | 人像修复1x  | model_id=36\n")
	sb.WriteString("  portrait_restore_2x  | 人像修复2x  | model_id=35\n\n")

	sb.WriteString("🏞️ 通用修复模型\n")
	sb.WriteString("───────────────────────────────────────\n")
	sb.WriteString("  general_restore_1x   | 通用修复1x  | model_id=31\n")
	sb.WriteString("  general_restore_2x   | 通用修复2x  | model_id=32\n")
	sb.WriteString("  general_restore_4x   | 通用修复4x  | model_id=33\n\n")

	sb.WriteString("🤖 高级模型（SD）\n")
	sb.WriteString("───────────────────────────────────────\n")
	sb.WriteString("  ultrahd_restore_2x   | 超高清修复2x | model_id=37\n")
	sb.WriteString("  generative_1x        | 生成式增强1x | model_id=38\n\n")

	sb.WriteString("═══════════════════════════════════════\n")
	sb.WriteString("📌 限制信息：\n")
	sb.WriteString("  最大输入: 36MP | 时长: 0.5秒~60分钟\n")
	sb.WriteString("  输入格式: .mp4 .mov .avi .mkv .webm .flv .ts 等\n")
	sb.WriteString("  输出格式: .mp4 .mov .mkv .m4v .avi .gif\n")
	sb.WriteString("  积分计算: 分辨率 × 帧率 × 时长\n")

	return protocol.SuccessResult(sb.String())
}
