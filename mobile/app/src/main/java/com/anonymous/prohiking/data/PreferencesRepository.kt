package com.anonymous.prohiking.data

import androidx.datastore.core.DataStore
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.intPreferencesKey
import androidx.datastore.preferences.core.longPreferencesKey
import androidx.datastore.preferences.core.stringPreferencesKey
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import java.sql.Time

interface PreferencesRepository {
    val userId: Flow<Int>
    val username: Flow<String>
    val password: Flow<String>

    val trailId: Flow<Int>
    val trailTime: Flow<Long>

    suspend fun updateUserId(id: Int)
    suspend fun updateUsernameAndPassword(username: String, password: String)
    suspend fun updateTrailId(id: Int)
    suspend fun updateTrailTime(time: Long)
}

class DefaultPreferencesRepository(private val dataStore: DataStore<Preferences>): PreferencesRepository {
    private companion object {
        val USER_ID = intPreferencesKey("user_id")
        val USERNAME = stringPreferencesKey("username")
        val PASSWORD = stringPreferencesKey("password")

        val TRAIL_ID = intPreferencesKey("trail_id")
        val TRAIL_TIME = longPreferencesKey("trail_time")
    }

    override val userId = dataStore.data.map { preferences ->
        preferences[USER_ID] ?: -1
    }

    override val username = dataStore.data.map { preferences ->
        preferences[USERNAME] ?: ""
    }

    override val password = dataStore.data.map { preferences ->
        preferences[PASSWORD] ?: ""
    }

    override val trailId = dataStore.data.map { preferences ->
        preferences[TRAIL_ID] ?: -1
    }

    override val trailTime = dataStore.data.map { preferences ->
        preferences[TRAIL_TIME] ?: -1
    }

    override suspend fun updateUserId(id: Int) {
        dataStore.edit { preferences ->
            preferences[USER_ID] = id
        }
    }

    override suspend fun updateUsernameAndPassword(username: String, password: String) {
        dataStore.edit { preferences ->
            preferences[USERNAME] = username
            preferences[PASSWORD] = password
        }
    }

    override suspend fun updateTrailId(id: Int) {
        dataStore.edit { preferences ->
            preferences[TRAIL_ID] = id
        }
    }

    override suspend fun updateTrailTime(time: Long) {
        dataStore.edit { preferences ->
            preferences[TRAIL_TIME] = time
        }
    }
}