package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"runs_adapter/adapter"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func NewHandler(w http.ResponseWriter, r *http.Request) {
	run := adapter.Run{}
	err := json.NewDecoder(r.Body).Decode(&run)
	if err != nil {
		WriteHTTPError(w, http.StatusInternalServerError, err)
		return
	}
	err = run.New()
	if err != nil {
		WriteHTTPError(w, http.StatusInternalServerError, err)
		return
	}
	w.Write([]byte("{\"Status\":\"OK\"}"))
}

func LastsHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	count, ok := params["count"]
	if ok {
		if strings.ToUpper(count) == "ALL" {
			runs, err := adapter.GetAllRuns()
			if err != nil {
				WriteHTTPError(w, http.StatusInternalServerError, err)
				return
			}
			err = json.NewEncoder(w).Encode(runs)
			if err != nil {
				WriteHTTPError(w, http.StatusInternalServerError, err)
				return
			}
			return
		} else {
			cnt, err := strconv.Atoi(count)
			if err != nil {
				WriteHTTPError(w, http.StatusInternalServerError, errors.New("Count is not \"ALL\" values and not integer type, "+err.Error()))
				return
			}
			runs, err := adapter.GetLastsRuns(cnt)
			if err != nil {
				WriteHTTPError(w, http.StatusInternalServerError, err)
				return
			}
			err = json.NewEncoder(w).Encode(runs)
			if err != nil {
				WriteHTTPError(w, http.StatusInternalServerError, err)
				return
			}
			return
		}
	}
	WriteHTTPError(w, http.StatusBadRequest, errors.New("params count not found"))
}

func RangeHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	start, startok := params["starttime"]
	if startok {
		end, endok := params["endttime"]
		if endok {
			s, err := strconv.Atoi(start)
			if err != nil {
				WriteHTTPError(w, http.StatusInternalServerError, errors.New("starttime is  not integer type, "+err.Error()))
				return
			}
			e, err := strconv.Atoi(end)
			if err != nil {
				WriteHTTPError(w, http.StatusInternalServerError, errors.New("endttime is  not integer type, "+err.Error()))
				return
			}
			runs, err := adapter.GetRangeRuns(s, e)
			if err != nil {
				WriteHTTPError(w, http.StatusInternalServerError, err)
				return
			}
			err = json.NewEncoder(w).Encode(runs)
			if err != nil {
				WriteHTTPError(w, http.StatusInternalServerError, err)
				return
			}
		}
	}
	WriteHTTPError(w, http.StatusBadRequest, errors.New("params starttime not found"))

}
