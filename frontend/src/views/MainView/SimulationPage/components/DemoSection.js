import {ChooseModel} from "./components/ChooseModel";
import {ChooseImage} from "./components/ChooseImage";
import React from "react";
import "bulma-extensions/bulma-divider/dist/css/bulma-divider.min.css"
import {runSimulation} from "../../../../API";
import {setInfoModalActive} from "../../../../components/InfoModal/InfoModal";
import {setErrorModalActive} from "../../../../components/ErrorModal/Error";
import {algorithmID} from "../../../utils";

const DemoSection = () => {
    const [disabledButton, setDisabledButton] = React.useState(true);

    const handleClick = async () => {
        const data = {
            model: document.getElementById('choose-model').value,
            image: document.getElementById('chosen-image').src,
            id: algorithmID
        }
        setDisabledButton(true);
        try {
            const resp = await runSimulation("demo", JSON.stringify(data))
            if (resp.status === 202) setInfoModalActive("Simulation started!")
        } catch (err) {
            setErrorModalActive(err)
        }
    }

    const setDisabled = () => {
        if (document.getElementById('chosen-image').src !== "https://bulma.io/images/placeholders/256x256.png"){
            setDisabledButton(false);
        } else {
            setDisabledButton(true);
        }
    }

    return (
        <div className="tile is-ancestor">
            <div className="tile is-parent">
                <div className="tile is-child box">
                    <ChooseModel/>
                    <br/>
                    <ChooseImage/>
                </div>
            </div>
            <div className="tile is-parent">
                <div className="tile is-child box">
                    <p className="subtitle">Chosen Image</p>
                    <img src="https://bulma.io/images/placeholders/256x256.png" id={"chosen-image"} alt={"chosen-image"} onLoad={setDisabled}/>
                </div>
            </div>
            <div className="tile is-parent is-vertical">
                <div className="tile is-child box">
                    <p className="subtitle">Run Simulation</p>
                    <div className="section has-text-centered">
                        <button id="run-simulation-button" disabled={disabledButton} className="button is-link" onClick={handleClick}>Run
                            Simulation
                        </button>
                    </div>
                </div>
            </div>
        </div>)
}
export {DemoSection}