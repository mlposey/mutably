import React, { Component } from 'react';
import SearchBar from './search-bar';
import InflectionTable from './inflection-table';
import ApiSearch from '../util/api-search.js';
import '../styles/app.css'

/**
 * The main application component
 * 
 * App encapsulates an input for looking up verbs and a table
 * for looking at categorized results.
 */
class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            search: new ApiSearch(),
            inflection: ApiSearch.getConjugationTable()
        };

        this.findAndStore('fiets');
    }

    /**
     * Performs an inflection call to the API and stores the result in
     * this.state.inflection
     * @param {string} verb the verb to look up
     */
    findAndStore(verb) {
        this.state.search.findVerb(verb, (inflection) => {
            this.setState({ inflection })
        });
    }

    render() {
        return (
            <div className="center-wrap">
                <div className="app">
                    <SearchBar onSearchWord={verb => this.findAndStore(verb)}/>
                    <InflectionTable inf={this.state.inflection} />
                </div>
            </div>
        );
    }
}

export default App;