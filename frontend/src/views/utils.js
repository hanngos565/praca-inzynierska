export const algorithmID = "alg1"
export const getExt = (filename) => filename.substr(filename.lastIndexOf('.'));
export const checkIfExists = (content, includes, errMessage) => {
    if (content.includes(includes)) throw new Error(errMessage);
};