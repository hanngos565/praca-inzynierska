import React from "react";
import {getImages, uploadImage} from "../../../../API";
import {ErrorMessage, setErrorModalActive} from "../../../../components/ErrorModal/Error";
import {getExt, checkIfExists} from "../../../utils";
import {InfoMessage, setInfoModalActive} from "../../../../components/InfoModal/InfoModal";

const extensions = [".png", ".jpg", ".jpeg"];
const convertToBase64 = (image) => {
    return new Promise((resolve) => {
        let reader = new FileReader();
        reader.readAsDataURL(image);

        reader.onload = (e) => {
            return resolve(reader.result)
        }
        reader.onerror = (e) => {
            console.log(e)
            throw new Error("Couldn't convert image to base64!");
        }
    })
};
const validateExt = (extension) => {
    if (!extensions.includes(extension)) throw new Error("Unsupported extension!");
};

const UploadImage = () => {
    const [contentImage, setContentImage] = React.useState();
    const [disabledButton, setDisabledButton] = React.useState(true);
    const [images, setImages] = React.useState();

    React.useEffect(() => {
        const fetchData = async () => {
            const i = await getImages();
            setImages(i.images);
        }
        fetchData()
            .catch(console.error);
    }, [contentImage])

    const clear = () => {
        setDisabledButton(true)
        setContentImage('')
        document.getElementById('file-input').value = null;
    }

    const setImage = async (e) => {
        try {
            validateExt(getExt(e.target.files[0].name));
            const base64Image = await convertToBase64(e.target.files[0]);
            if (images) {
                checkIfExists(images, base64Image, "This image is already in the database!")
            }
            setContentImage(base64Image);
            setDisabledButton(false);
        } catch (err) {
            clear()
            setErrorModalActive(err)
        }
    }

    const handleClick = async () => {
        try {
            const jsonData = {
                content: contentImage.toString(),
            };
            let resp = await uploadImage(jsonData)
            clear()
            setImages('')
            if (resp.status === 200) setInfoModalActive("Image uploaded!")
        }catch(err){
            console.log(err)
        }
    }

    return (
        <div className="container has-text-centered">
            <h5>Accepted image formats: {extensions}</h5><br/>
            <div className="file is-centered">
                <label className="file-label">
                    <input className="file-input" id="file-input" type="file" name="resume" accept={extensions}
                           onChange={(e) => setImage(e)}/>
                    <span className="file-cta">
                        <span className="file-icon">
                            <i className="fas fa-upload"/>
                        </span>
                        <span className="file-label">
                            Choose an image…
                        </span>
                    </span>
                </label>
            </div>
            <br/>
            <button className="button is-primary" disabled={disabledButton} onClick={handleClick}>Upload Image!</button>
            <ErrorMessage/>
            <InfoMessage/>
        </div>
    )
}

export {UploadImage}