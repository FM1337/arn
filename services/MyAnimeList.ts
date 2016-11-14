import * as xml2js from 'xml2js'
import * as arn from 'arn'
import * as request from 'request-promise'
import { Anime } from 'arn/interfaces/Anime'

const COMPLETED = 2
const WATCHING = 1

interface MyAnimeListAnime {
	id: number
	similarity: number
	my_status: string
	my_watched_episodes: string
	series_animedb_id: string
	series_episodes: string
	series_title: string
	series_image: string
}

class MyAnimeList {
	static headers = {
		'User-Agent': 'Anime Release Notifier'
	}

	static xmlParser = new xml2js.Parser({
		explicitArray: false,	// Don't put single nodes into an array
		ignoreAttrs: true,		// Ignore attributes and only create text nodes
		trim: true,
		normalize: true,
		explicitRoot: false
	})

	getAnimeListUrl(userName) {
		return `http://myanimelist.net/animelist/${userName}&status=${WATCHING}`
	}

	getAnimeList(userName, callback) {
		return request({
			uri: `http://myanimelist.net/malappinfo.php?u=${userName}&status=all&type=anime`,
			method: 'GET',
			headers: MyAnimeList.headers
		}).then(body => {
			MyAnimeList.xmlParser.parseString(body, (error, json) => {
				if(error) {
					callback(error, [])
					return
				}

				let watching: any[] = []

				if(!Array.isArray(json.anime)) {
					return callback(undefined, watching)
				}

				let malAnime: MyAnimeListAnime[] = json.anime

				watching = malAnime.filter(entry => parseInt(entry.my_status) === WATCHING)
				.map(entry => {
					let episodesWatched = parseInt(entry.my_watched_episodes)
					let nextEpisodeToWatch = episodesWatched + 1
					let episodesOffset = 0

					return {
						title: entry.series_title,
						image: entry.series_image,
						url: 'http://myanimelist.net/anime/' + entry.series_animedb_id,
						providerId: parseInt(entry.series_animedb_id),
						airingDate: null,
						episodes: {
							watched: episodesWatched ? episodesWatched : 0,
							next: nextEpisodeToWatch,
							available: 0,
							max: entry.series_episodes ? parseInt(entry.series_episodes) : -1,
							offset: episodesOffset
						}
					}
				})

				let tasks: Promise<any>[] = []
				watching.forEach(anime => {
					tasks.push(arn.getAnimeIdBySimilarTitle(anime, 'MyAnimeList').then(match => {
						anime.id = match ? match.id : null
						anime.similarity = match ? match.similarity : 0
					}))
				})

				Promise.all(tasks).then(() => callback(undefined, watching))
			})
		}).catch(error => {
			callback(error, [])
		})
	}

	getUserImage(userName) {
		return request({
			uri: `http://myanimelist.net/malappinfo.php?u=${userName}`,
			method: 'GET',
			headers: MyAnimeList.headers
		}).then(body => {
			let idRegex = /<user_id>(\d+)<\/user_id>/
			let match = body.match(idRegex)

			if(match && match[1])
				return `https://myanimelist.cdn-dena.com/images/userimages/${match[1]}.jpg`
			else
				throw new Error('Couldn\'t get MAL user ID')
		})
	}
}

module.exports = new MyAnimeList()