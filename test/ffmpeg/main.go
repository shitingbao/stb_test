package main

// ffmpeg -i "https://customer-6qs66h96fgl8gzoy.cloudflarestream.com/1b8dcf9b335445bab91c848d90c4420d/manifest/video.m3u8" -c copy -bsf:a aac_adtstoasc output_video.mp4
// ffmpeg -i input.mp4 -ss 2 -c copy output.mp4
// 解释
// -i input.mp4：指定输入视频文件。
// -ss 2：设置开始时间为2秒，ffmpeg将从视频的第2秒开始剪辑。可以小数点比如0.5秒
// -c copy：复制音频和视频流，而不重新编码。这样操作会更快并保持原始视频质量。
// output.mp4：指定输出文件名。
// 这个命令将会跳过视频的前两秒，并保留其余部分。

// https://customer-6qs66h96fgl8gzoy.cloudflarestream.com/1b8dcf9b335445bab91c848d90c4420d/manifest/video.m3u8

// 获取一个视频里面的一段时间的内容
// ./ffmpeg -i input_video.mp4 -ss 00:00:05 -t 00:00:05 -c copy output_video.mp4
// -ss 00:00:05：表示从 5 秒钟开始截取。
// -t 00:00:05：表示截取 5 秒的长度。
// ./ffmpeg -i aa.mp4 -ss 00:00:01 -t 00:00:05 output_video.mp4
