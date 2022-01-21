import {setImagesModalActive} from "../ImagesModal";
import React from "react";

const ChooseImage = () => {
    return (
        <div>
            <p className="subtitle">Choose Image</p>
            <div className="container has-text-centered">
                <button className="button is-link" onClick={setImagesModalActive}>Choose Image!</button>
            </div>
        </div>
    )
}

export {ChooseImage}