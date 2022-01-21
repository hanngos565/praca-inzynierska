export const setErrorModalActive = (err) => {
    const modal = document.getElementById("error-modal")
    modal.classList.add('is-active')
    document.getElementById("error-message").innerText = err
}

const unsetErrorModalActive = () => {
    const modal = document.getElementById("error-modal")
    modal.classList.remove('is-active')
}


const ErrorMessage = () => {
    return (
        <div className="modal" id="error-modal">
            <div className="modal-background"/>
            <div className="modal-card has-background-danger">
                <button className="delete" aria-label="close" onClick={unsetErrorModalActive}/>
                <section className="modal-card-body">
                    <div id="error-message"/>
                </section>
            </div>
        </div>)

}

export {ErrorMessage}