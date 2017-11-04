/** Facilitates word searches using the Mutably API */
class ApiSearch {
    constructor() {
        this.apiV1 = 'http://srv.marcusposey.com:9000/api/v1';
    }

    /**
     * Retrieves the conjugation table for a verb
     * A valid but empty table is passed to callback if verb does not exist
     * @param {string} verb the verb to find
     * @param {Function} callback a function which takes an object (the table)
     *                            as an argument
     */
    findVerb(verb, callback) {
        const resource = this.apiV1 + '/words/' + verb + '/inflections';
        fetch(resource)
            .then(response => {
                if (response.status === 404) {
                    return ApiSearch.getConjugationTable();
                }
                return response.json();
            })
            .then(callback);
    }

    /** Retrieves a blank conjugation table */
    static getConjugationTable() {
        return {
            Infinitive: '',
            Present: {First: [], Second: [], Third: [], Plural: []},
            Past: {First: [], Second: [], Third: [], Plural: []}
        }
    }
}

export default ApiSearch;