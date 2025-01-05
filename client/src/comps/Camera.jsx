import React, { useState, useRef, useEffect } from "react";
import {Camera} from "react-camera-pro";
import { API_URL } from "./socket";

const Cam = ({setImgUUID}) => {
  const camera = useRef(null);

  const takePhoto = () => {
    fetch(`${API_URL}/api/upload`, {
      method: "POST",
      body: camera.current.takePhoto(),
      headers: {
        "Content-type": "text/plain"
      }
    }).then((r) => {
      r.text().then((uid) => setImgUUID(uid))
    });
  }

  // useEffect(() => {
  //   if (image != null) {
  //     fetch(`http://localhost:80/api/image/get/${image}`, {
  //       method: "GET"
  //     }).then((r) => {
  //       r.text().then((dat) => setImageData(dat))
  //     })
  //   }
  // }, [image])

  return (
    <div className="Camera">
      <Camera ref={camera} aspectRatio={1} />
      <button style={{position: "absolute", top:"0", left:"0", zIndex: "45"}} onClick={takePhoto}>Take photo</button>
    </div>
  );
}

export default Cam;