export const setInfoModalActive = (msg) => {
    const modal = document.getElementById("info-modal")
    modal.classList.add('is-active')
    document.getElementById("info-message").innerText = msg
}

const unsetInfoModalActive = () => {
    const modal = document.getElementById("info-modal")
    modal.classList.remove('is-active')
    window.location.reload();
}


const InfoMessage = () => {
    return (
        <div className="modal" id="info-modal">
            <div className="modal-background"/>
            <div className="modal-card has-background-success">
                <button className="delete" aria-label="close" onClick={unsetInfoModalActive}/>
                <section className="modal-card-body">
                    <div id="info-message"/>
                </section>
            </div>
        </div>)

}

export {InfoMessage}