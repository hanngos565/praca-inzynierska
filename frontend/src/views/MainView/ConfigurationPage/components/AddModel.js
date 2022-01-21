import React from "react";
import {getModels, addModel} from "../../../../API";
import {getExt, checkIfExists} from "../../../utils";
import {ErrorMessage, setErrorModalActive} from "../../../../components/ErrorModal/Error";
import {setInfoModalActive} from "../../../../components/InfoModal/InfoModal";
import {algorithmID} from "../../../utils";

export const algorithm = {"alg1": [".h5"], "alg2": ".caffe", "1": [".h5", ".png"]};


const validateExt = (alg, extension) => {
    let extensions = [];
    for (let key in algorithm)
        if (key === alg) extensions = algorithm[key];
    if (!extensions.includes(extension)) throw new Error('Unsupported extension!');
};

const AddModel = () => {

    const [uploadedModel, setUploadedModel] = React.useState();
    let [modelName, setModelName] = React.useState();
    let [disabledButton, setDisabledButton] = React.useState(true);
    const [models, setModels] = React.useState();
    const [disableCustomName, setDisableCustomName] = React.useState(true);
    const [ext, setExt] = React.useState();

    React.useEffect(() => {
        const fetchData = async () => {
            const m = await getModels(algorithmID);
            setModels(m);
        }
        fetchData()
            .catch(console.error);
    }, [])

    const setModel = (e) => {
        clear();
        try {
            const extension = getExt(e.target.files[0].name);
            validateExt(e.target.id, extension);
            setExt(extension);
            setUploadedModel(e.target.files[0]);
            document.getElementById("file-name").innerText = e.target.files[0].name;
            checkIfExists(models, e.target.files[0].name, "Model with such a name already exists!\n" +
                "Please provide new name!");
            setModelName(e.target.files[0].name);
            setDisabledButton(false);
        } catch (err) {
            if (err.toString() === "Error: Model with such a name already exists!\n" +
            "Please provide new name!") {
                setDisableCustomName(false);
            } else if (err.toString() === 'Error: Unsupported extension!') {
                e.target.value = null;
            }
            setErrorModalActive(err)
        }
    }

    const setName = async (e) => {
        try {
            if(e.target.value === '') {
                setDisabledButton(true)
            } else {
                checkIfExists(models, e.target.value + ext, "Model with such a name already exists!\n" +
                    "Please provide new name!");
                setModelName(e.target.value + ext)
                setDisabledButton(false)
            }
        } catch (err) {
            setDisabledButton(true)
            setErrorModalActive(err)
        }
    }

    const clear = () => {
        setExt('');
        setDisableCustomName(true);
        setUploadedModel('');
        setModelName('');
        setDisabledButton(true);
        document.getElementById("custom-name").value = null;
        document.getElementById("file-name").innerText = "File name"
    }

    const handleClick = async (e) => {
        setDisabledButton(true)
        const data = new FormData()
        data.append('id', e.target.id)
        data.append('name', modelName)
        data.append('model', new Blob([uploadedModel], {type:"application/octet-stream"}))

        try {
            const resp = await addModel(data)
            if (resp.status === 200) setInfoModalActive("Model added!")
        } catch (err) {
            setErrorModalActive(err)
        }
    }

    return (
        <div className="container has-text-centered">
            <h5>Accepted model formats: {algorithm.alg1}</h5><br/>
            <div className="file is-centered has-name is-fullwidth">
                <label className="file-label">
                    <input className="file-input" id={algorithmID} type="file" name="resume" accept={algorithm.alg1}
                           onChange={(e) => setModel(e)}/>
                    <span className="file-cta">
                        <span className="file-icon">
                            <i className="fas fa-upload"/>
                        </span>
                        <span className="file-label">
                            Choose a model…
                        </span>
                    </span>
                    <span className="file-name" id="file-name">
                        File name
                    </span>
                </label>
            </div>
            <br/>
            <input className="input" disabled={disableCustomName} id="custom-name" type="text" placeholder="New name"
                   onChange={e => setName(e)}/><br/>
            <br/>
            <button id={algorithmID} className="button is-primary" disabled={disabledButton} onClick={handleClick}>Add
                model!
            </button>
            <ErrorMessage/>
        </div>
    )
}

export {AddModel}