# rt2togpx

Converts rt2 route files(Some oziexplorer format) to gpx route files.

## RT2

    W,RTP315,61.313424,26.696061,0
    W,RTP318,61.313770,26.695530,0
    W,RTP321,61.314430,26.695110,0
    W,RTP326,61.315230,26.693790,0

## GPX
    
    <rtept lat="61.313424" lon="26.696061"></rtept>
    <rtept lat="61.31377" lon="26.69553"></rtept>
    <rtept lat="61.31443" lon="26.69511"></rtept>
    <rtept lat="61.31523" lon="26.69379"></rtept>

## Usage

    $ ./rt2togpx -help
    Usage of ./rt2togpx:
    -out string
      output gpx file name (default "route.gpx")
    -rt2file string
      rt2 file to be converted
    -verbose
      Print some details about route

e.g.

     $ ./rt2togpx -rt2file route.rt2 
     2018/06/29 18:43:56 GPX file succesfully saved to route.gpx
