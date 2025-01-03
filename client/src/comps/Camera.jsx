import React, { useState, useRef } from "react";
import {Camera} from "react-camera-pro";

const Component = () => {
  const camera = useRef(null);
  const [image, setImage] = useState(null);

  return (
    <div>
      <Camera ref={camera} aspectRatio="4/3" />
      <button style={{position: "absolute", top:"0", left:"0", zIndex: "45"}} onClick={() => setImage(camera.current.takePhoto())}>Take photo</button>
      <img src={image} alt='Taken photo'/>
    </div>
  );
}

export default Component;