import {UploadImage} from "./components/UploadImage";
import React from "react";
import {AddModel} from "./components/AddModel";

const Configuration = () => {
    return (
        <section id={"configuration"} className="section notification is-primary">
            <h1 className="title">Configuration</h1>
            <div className="tile is-ancestor">
                <div className="tile is-parent">
                    <div className="tile is-child box">
                        <p className="subtitle">Upload Image</p>
                        <UploadImage/>
                    </div>
                </div>
                <div className="tile is-parent">
                    <div className="tile is-child box">
                        <p className="subtitle">Add Model</p>
                        <AddModel/>
                    </div>
                </div>
            </div>
        </section>
    )
}

export {Configuration}