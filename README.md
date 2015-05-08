## Usage

15:07 $ vostok-transcoder-tester -h
Usage of vostok-transcoder-tester:
  -file="": name of the file in '/incoming' that you want to transcode
  -idlist="": comma-separated list of dictionary identifiers to 'start'
  -url="https://reeldx-vostok-tc-stag.elasticbeanstalk.com/api/v1/transcode/start": URL for starting transcode jobs on vostok-server
  -v=false: prints the constructed curl command

15:07 $ vostok-transcoder-tester -file=c2.mp4 -idlist=0240,0360 -v

curl -k -H 'Content-type: application/json' -d '{"transcode_videos":[{"file_name_encoded":"c2.mp4","transcode_entry":{"audio_bit_rate":"125k","audio_codec":"libfdk_aac","audio_sample_rate":"44100","create":"true","horizontal_size":"640","identifier":"0360","key_frame_interval":"1","niceness":4,"output_file_extension":".mp4","revision":"-","vertical_size":"360","video_bit_rate":"400k","video_codec":"libx264","view":"true"}},{"file_name_encoded":"c2.mp4","transcode_entry":{"audio_bit_rate":"64k","audio_codec":"libmp3lame","audio_sample_rate":"22050","create":"true","horizontal_size":"426","identifier":"0240","key_frame_interval":"1","niceness":0,"output_file_extension":".mp4","revision":"-","vertical_size":"240","video_bit_rate":"250k","video_codec":"flv","view":"false"}}]}' -i https://reeldx-vostok-tc-stag.elasticbeanstalk.com/api/v1/transcode/start
HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Content-Type: application/json; charset=utf-8
Date: Fri, 08 May 2015 22:08:45 GMT
Server: nginx/1.6.2
Strict-Transport-Security: max-age=315360000
X-Frame-Options: DENY
X-Xss-Protection: 1; mode=block
Content-Length: 817
Connection: keep-alive

{"transcode_videos":[{"added":"2015-05-08T22:08:45.810777047Z","file_name_encoded":"c2.mp4","transcode_entry":{"identifier":"0360","revision":"-","niceness":4,"create":"true","view":"true","key_frame_interval":"1","video_codec":"libx264","video_bit_rate":"400k","horizontal_size":"640","vertical_size":"360","audio_codec":"libfdk_aac","audio_bit_rate":"125k","audio_sample_rate":"44100","output_file_extension":".mp4"}},{"added":"2015-05-08T22:08:45.810777047Z","file_name_encoded":"c2.mp4","transcode_entry":{"identifier":"0240","revision":"-","niceness":0,"create":"true","view":"false","key_frame_interval":"1","video_codec":"flv","video_bit_rate":"250k","horizontal_size":"426","vertical_size":"240","audio_codec":"libmp3lame","audio_bit_rate":"64k","audio_sample_rate":"22050","output_file_extension":".mp4"}}]}

