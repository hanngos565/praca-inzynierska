import React from "react";
import {getImages} from "../../../../API";

export const setImagesModalActive = () => {
    const modal = document.getElementById("modal")
    modal.classList.add('is-active')
}

const ImagesModal = () => {
    const [images, setImages] = React.useState();

    const unsetImagesModalActive = () => {
        const modal = document.getElementById("modal")
        modal.classList.remove('is-active')
    }
    const chooseImage = (e) => {
        document.getElementById('chosen-image').src = e.target.src
        unsetImagesModalActive()
    }

    React.useEffect(() => {
        const fetchData = async () => {
            const i = await getImages();
            setImages(i.images);
        }
        fetchData()
            .catch(console.error);
    }, [])


    return (
        <div id="modal" className="modal">
            <div className="modal-background"/>
            <div className="modal-card">
                <header className="modal-card-head">
                    <p className="modal-card-title">Choose Image</p>
                    <button className="delete" aria-label="close" onClick={unsetImagesModalActive}/>
                </header>
                <section className="modal-card-body has-text-centered">
                    {images &&
                        images.map(i => {
                            return (
                                <figure key={i} className="image is-inline-block">
                                    <img src={i} alt={i} onClick={e => chooseImage(e)}/>
                                </figure>)
                        })}
                </section>
            </div>
        </div>)
}
export {ImagesModal}
